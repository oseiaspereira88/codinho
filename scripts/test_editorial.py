#!/usr/bin/env python3
"""Deterministic tests using isolated Git repositories and a synthetic author."""
import argparse
import contextlib
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import threading
import time
import unittest
from unittest.mock import patch

SPEC = importlib.util.spec_from_file_location('editorial', Path(__file__).with_name('editorial.py'))
editorial = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(editorial)

FAKE = '''#!/usr/bin/env python3
import json, pathlib, sys, time
prompt = sys.stdin.read()
pathlib.Path("a.txt").write_text("revised\\n" if "resume" in sys.argv else "edited\\n")
print(json.dumps({"type":"thread.started", "thread_id":"12345678-1234-1234-1234-123456789abc"}), flush=True)
if "OUTSIDE" in prompt: pathlib.Path("outside.txt").write_text("bad")
if "SLOW" in prompt: time.sleep(10)
'''


class EditorialTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.repo = self.root / 'repo'
        self.repo.mkdir()
        self.run = self.root / 'run'
        self.queue = self.root / 'queue.json'
        self.fake = self.root / 'codex'
        self.fake.write_text(FAKE)
        self.fake.chmod(0o755)
        subprocess.run(['git', 'init', '-q', str(self.repo)], check=True)
        editorial.git(self.repo, 'config', 'user.email', 'test@example.invalid')
        editorial.git(self.repo, 'config', 'user.name', 'Test')
        (self.repo / 'a.txt').write_text('initial\n')
        editorial.git(self.repo, 'add', '.')
        editorial.git(self.repo, 'commit', '-qm', 'initial')

    def initialize(self, **overrides):
        tasks = overrides.pop('tasks', [dict(id='a', title='Edit', paths=['a.txt'], acceptance=['meaningful content'])])
        self.queue.write_text(json.dumps(dict(spec='example', tasks=tasks)))
        values = dict(run_dir=str(self.run), queue=str(self.queue), parallelism=2,
                      batch_size=1, timeout_seconds=10, model='gpt-5.6-luna',
                      reasoning='high', codex=str(self.fake), completed_task=[], task=[])
        values.update(overrides)
        with contextlib.chdir(self.repo):
            return editorial.init(argparse.Namespace(**values))

    def batch(self, name='batch-001'):
        return self.run / name

    def test_invalid_configuration_and_queue(self):
        for overrides in ({'parallelism': 9}, {'parallelism': 0}, {'batch_size': 0},
                          {'timeout_seconds': 0}, {'tasks': [dict(id='bad', title='Bad', paths=['../escape'], acceptance=['x'])]}):
            with self.subTest(overrides=overrides), self.assertRaises(ValueError):
                self.initialize(**overrides)
        self.initialize()
        with self.assertRaises(ValueError):
            self.initialize()

    def test_run_review_resume_integrate(self):
        config = self.initialize()
        state = editorial.execute(config, self.batch())
        self.assertEqual(state['status'], 'awaiting-review')
        with self.assertRaisesRegex(ValueError, 'approval'):
            editorial.integrate(config, self.batch())
        editorial.review(config, self.batch(), 'request-changes', 'Clarify evidence')
        original = subprocess.Popen
        commands = []

        def capture(command, **kwargs):
            commands.append(command)
            return original(command, **kwargs)

        with patch.object(editorial.subprocess, 'Popen', side_effect=capture):
            state = editorial.execute(config, self.batch(), 'Clarify evidence')
        self.assertEqual(state['status'], 'awaiting-review')
        resumed = [c for c in commands if 'resume' in c][0]
        self.assertEqual(resumed[3], '12345678-1234-1234-1234-123456789abc')
        self.assertNotIn('--last', resumed)
        editorial.review(config, self.batch(), 'approve', 'Checks and content reviewed')
        result = editorial.integrate(config, self.batch())
        self.assertEqual(result['status'], 'integrated')
        self.assertEqual((self.repo / 'a.txt').read_text(), 'revised\n')
        self.assertIn(b'POSE-Spec: example', editorial.git(self.repo, 'log', '-1', '--format=%B'))

    def test_out_of_scope_and_stale_approval(self):
        config = self.initialize(tasks=[dict(id='a', title='OUTSIDE', paths=['a.txt'], acceptance=['x'])])
        result = editorial.execute(config, self.batch())
        self.assertEqual(result['status'], 'failed')
        self.assertIn('out-of-scope', result['error'])
        self.assertTrue((self.batch() / 'worktree/outside.txt').exists())

    def test_stale_and_dirty_target(self):
        config = self.initialize()
        editorial.execute(config, self.batch())
        editorial.review(config, self.batch(), 'approve', 'Reviewed')
        (self.repo / 'a.txt').write_text('user edit\n')
        with self.assertRaisesRegex(ValueError, 'clean'):
            editorial.integrate(config, self.batch())
        (self.batch() / 'worktree/a.txt').write_text('post-review edit\n')
        with self.assertRaisesRegex(ValueError, 'stale'):
            editorial.integrate(config, self.batch())
        self.assertEqual((self.repo / 'a.txt').read_text(), 'user edit\n')

    def test_parallel_run_isolated_batches(self):
        tasks = [dict(id=str(i), title='Edit', paths=[f'{i}.txt'], acceptance=['x']) for i in range(4)]
        self.initialize(tasks=tasks)
        self.fake.write_text(FAKE.replace('"a.txt"', 'json.loads(prompt.splitlines()[1])[0]["paths"][0]'))
        active = peak = 0
        lock = threading.Lock()
        original = editorial.execute

        def counted(config, folder):
            nonlocal active, peak
            with lock:
                active += 1
                peak = max(peak, active)
            try:
                time.sleep(0.05)
                return original(config, folder)
            finally:
                with lock:
                    active -= 1

        with contextlib.redirect_stdout(None), patch.object(editorial, 'execute', side_effect=counted):
            code = editorial.main(['run', '--run-dir', str(self.run)])
        self.assertEqual(code, 0)
        self.assertEqual(peak, 2)
        for i in range(1, 5):
            state = editorial.read(self.batch(f'batch-{i:03d}') / 'state.json')
            self.assertEqual(state['status'], 'awaiting-review' if i <= 2 else 'pending')
            if i <= 2:
                self.assertEqual((self.batch(f'batch-{i:03d}') / f'worktree/{i-1}.txt').read_text(), 'edited\n')
        self.assertEqual((self.repo / 'a.txt').read_text(), 'initial\n')

    def test_timeout_preserves_work(self):
        config = self.initialize(timeout_seconds=1, tasks=[dict(id='a', title='SLOW', paths=['a.txt'], acceptance=['x'])])
        result = editorial.execute(config, self.batch())
        self.assertEqual(result['status'], 'failed')
        self.assertIn('timed out', result['error'])
        self.assertEqual(result['session_id'], '12345678-1234-1234-1234-123456789abc')
        self.assertEqual((self.batch() / 'worktree/a.txt').read_text(), 'edited\n')

    def test_pilot_limits_dispatch_and_uncreated_allowed_path(self):
        tasks = [dict(id=str(i), title='Edit', paths=['a.txt', 'future.txt'], acceptance=['x']) for i in range(3)]
        self.initialize(tasks=tasks)
        with contextlib.redirect_stdout(None):
            self.assertEqual(editorial.main(['run', '--run-dir', str(self.run), '--limit-batches', '1']), 0)
        self.assertEqual(editorial.read(self.batch() / 'state.json')['status'], 'awaiting-review')
        self.assertEqual(editorial.read(self.batch('batch-002') / 'state.json')['status'], 'pending')

    def test_hidden_staged_path_is_rejected(self):
        config = self.initialize()
        editorial.execute(config, self.batch())
        work = self.batch() / 'worktree'
        (work / 'outside.txt').write_text('hidden change')
        editorial.git(work, 'add', 'outside.txt')
        (work / 'outside.txt').unlink()
        with self.assertRaisesRegex(ValueError, 'out-of-scope'):
            editorial.review(config, self.batch(), 'approve', 'Reviewed')

    def test_retry_before_worktree_and_session_exist(self):
        config = self.initialize()
        with patch.object(editorial, 'git', side_effect=OSError('worktree unavailable')):
            self.assertEqual(editorial.execute(config, self.batch())['status'], 'failed')
        self.assertEqual(editorial.execute(config, self.batch(), 'Retry')['status'], 'awaiting-review')

    def test_no_change_audit_and_empty_evidence(self):
        config = self.initialize()
        self.fake.write_text(FAKE.replace('pathlib.Path("a.txt").write_text("revised\\n" if "resume" in sys.argv else "edited\\n")', 'pass'))
        editorial.execute(config, self.batch())
        with self.assertRaisesRegex(ValueError, 'nonempty'):
            editorial.review(config, self.batch(), 'approve', ' ')
        editorial.review(config, self.batch(), 'approve', 'Audit confirms existing content')
        self.assertTrue(editorial.integrate(config, self.batch())['no_changes'])

    def test_dependencies_wait_for_integration_and_use_new_base(self):
        tasks = [dict(id='first', title='Edit', paths=['a.txt'], acceptance=['x']),
                 dict(id='second', title='Edit', paths=['a.txt'], acceptance=['x'], depends_on=['first'])]
        config = self.initialize(tasks=tasks, batch_size=5)
        folders = sorted(self.run.glob('batch-*/state.json'))
        self.assertEqual(editorial.ready_batches(config, folders, 8), [self.batch()])
        editorial.execute(config, self.batch())
        self.assertEqual(editorial.ready_batches(config, folders, 8), [])
        editorial.review(config, self.batch(), 'approve', 'Checks passed')
        first = editorial.integrate(config, self.batch())
        self.assertEqual(editorial.ready_batches(config, folders, 8), [self.batch('batch-002')])
        state = editorial.execute(config, self.batch('batch-002'))
        self.assertEqual(state['base'], first['commit'])

    def test_completed_task_unblocks_first_remaining_batches(self):
        tasks = [dict(id='pilot', title='Pilot', paths=['pilot.txt'], acceptance=['x']),
                 dict(id='first', title='First', paths=['a.txt'], acceptance=['x'], depends_on=['pilot'])]
        config = self.initialize(tasks=tasks, completed_task=['pilot'])
        self.assertEqual(editorial.read(self.batch() / 'state.json')['tasks'][0]['id'], 'first')
        self.assertEqual(editorial.ready_batches(config, sorted(self.run.glob('batch-*/state.json')), 1), [self.batch()])

    def test_selected_tasks_keep_requested_order(self):
        tasks = [dict(id='one', title='One', paths=['one.txt'], acceptance=['x']),
                 dict(id='two', title='Two', paths=['two.txt'], acceptance=['x'])]
        self.initialize(tasks=tasks, task=['two'])
        self.assertEqual(editorial.read(self.batch() / 'state.json')['tasks'][0]['id'], 'two')

    def test_approved_delivery_can_be_revised(self):
        config = self.initialize()
        editorial.execute(config, self.batch())
        editorial.review(config, self.batch(), 'approve', 'Checked')
        state = editorial.execute(config, self.batch(), 'Reconsider approval')
        self.assertEqual(state['status'], 'awaiting-review')
        self.assertNotIn('approved_digest', state)

    def test_recover_rejects_live_author_and_restores_explicit_session(self):
        config = self.initialize()
        state = editorial.execute(config, self.batch())
        state.update(status='running', author_pid=os.getpid())
        state.pop('session_id')
        editorial.save(self.batch() / 'state.json', state)
        with self.assertRaisesRegex(ValueError, 'still exists'):
            editorial.recover(config, self.batch())
        with patch.object(editorial.os, 'kill', side_effect=ProcessLookupError):
            result = editorial.recover(config, self.batch())
        self.assertEqual(result['status'], 'failed')
        self.assertEqual(result['session_id'], '12345678-1234-1234-1234-123456789abc')

    def test_reconcile_failed_commit(self):
        (self.repo / 'a.txt').write_text(''.join(f'line {i}\n' for i in range(30)))
        editorial.git(self.repo, 'commit', '-am', 'Long document')
        self.fake.write_text(FAKE.replace('"edited\\n"', 'pathlib.Path("a.txt").read_text().replace("line 0\\n", "edited\\n")'))
        config = self.initialize()
        editorial.execute(config, self.batch())
        editorial.review(config, self.batch(), 'approve', 'Checked')
        (self.repo / 'a.txt').write_text((self.repo / 'a.txt').read_text().replace('line 28\n', 'independent\n'))
        editorial.git(self.repo, 'commit', '-am', 'Independent edit')
        hook = self.repo / '.git/hooks/pre-commit'
        hook.write_text('#!/bin/sh\nexit 1\n')
        hook.chmod(0o755)
        with self.assertRaisesRegex(ValueError, 'commit failed'):
            editorial.integrate(config, self.batch())
        with self.assertRaisesRegex(ValueError, 'clean target'):
            editorial.reconcile(config, self.batch())
        hook.unlink()
        editorial.git(self.repo, 'commit', '-m', 'Manual integration', '-m', 'POSE-Spec: example')
        self.assertEqual(editorial.reconcile(config, self.batch())['status'], 'integrated')


if __name__ == '__main__':
    unittest.main()
