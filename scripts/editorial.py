#!/usr/bin/env python3
"""Persistent, manually reviewed Codex editorial batches (Python stdlib only)."""
import argparse
import concurrent.futures
import fcntl
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import signal
import subprocess
import sys
import tempfile
import uuid


def read(path):
    return json.loads(Path(path).read_text())


def save(path, value):
    path = Path(path)
    temp = path.with_suffix('.tmp')
    temp.write_text(json.dumps(value, ensure_ascii=False, indent=2) + '\n')
    temp.replace(path)


def git(repo, *args):
    return subprocess.check_output(['git', '-C', str(repo), *args], stderr=subprocess.PIPE)


def clean(repo):
    return not git(repo, 'status', '--porcelain').strip()


def valid_path(path):
    return (isinstance(path, str) and bool(path) and
            not PurePosixPath(path).is_absolute() and
            all(part not in ('', '.', '..', '.git') for part in path.split('/')) and
            not path.startswith(':') and
            not any(c in path for c in '\n\r\0\\*?['))


def init(args):
    root = Path(args.run_dir).resolve()
    if root.exists():
        raise ValueError('run-dir must not exist; configuration is immutable')
    repo = Path(git(Path.cwd(), 'rev-parse', '--show-toplevel').decode().strip())
    if root == repo or repo in root.parents:
        raise ValueError('run-dir must be outside the repository')
    queue = read(args.queue)
    if not isinstance(queue, dict):
        raise ValueError('queue must be an object')
    if not re.fullmatch(r'[a-z0-9][a-z0-9-]*', queue.get('spec', '')):
        raise ValueError('queue requires a valid POSE spec slug')
    source_tasks = queue.get('tasks')
    if not isinstance(source_tasks, list) or not source_tasks:
        raise ValueError('queue requires tasks')
    source_ids = {task.get('id') for task in source_tasks if isinstance(task, dict)}
    completed = set(args.completed_task)
    if not completed <= source_ids:
        raise ValueError('completed-task must identify queue tasks')
    selected = args.task or [task.get('id') for task in source_tasks]
    if len(selected) != len(set(selected)) or not set(selected) <= source_ids:
        raise ValueError('task must identify unique queue tasks')
    by_id = {task['id']: task for task in source_tasks}
    tasks = [by_id[task_id] for task_id in selected if task_id not in completed]
    if not tasks:
        raise ValueError('queue has no remaining tasks')
    ids = set()
    for task in tasks:
        if not isinstance(task, dict) or not isinstance(task.get('id'), str) or not task['id'] or task['id'] in ids:
            raise ValueError('task IDs must be nonempty and unique')
        deps = task.get('depends_on', [])
        if not isinstance(deps, list) or not all(isinstance(d, str) and (d in ids or d in completed) for d in deps):
            raise ValueError('depends_on must reference earlier task IDs')
        ids.add(task['id'])
        if not isinstance(task.get('title'), str) or not task['title'].strip():
            raise ValueError('task requires title')
        if not isinstance(task.get('paths'), list) or not task['paths'] or not all(valid_path(p) for p in task['paths']):
            raise ValueError('paths must be exact safe repository-relative files')
        if not isinstance(task.get('acceptance'), list) or not task['acceptance'] or not all(isinstance(a, str) and a.strip() for a in task['acceptance']):
            raise ValueError('task requires acceptance criteria')
    if not 1 <= args.parallelism <= 8 or args.batch_size < 1 or args.timeout_seconds < 1:
        raise ValueError('parallelism must be 1..8; batch-size and timeout must be positive')
    if not args.model.strip() or args.model.startswith('-'):
        raise ValueError('model must be nonempty')
    config = dict(repo=str(repo), base=git(repo, 'rev-parse', 'HEAD').decode().strip(),
                  queue=queue, model=args.model, reasoning=args.reasoning,
                  parallelism=args.parallelism, batch_size=args.batch_size,
                  timeout_seconds=args.timeout_seconds, codex=args.codex,
                  completed_task_ids=sorted(completed))
    root.mkdir(parents=True)
    save(root / 'config.json', config)
    batches = []
    current = []
    for task in tasks:
        if current and (len(current) == args.batch_size or
                        set(task.get('depends_on', [])) & {t['id'] for t in current}):
            batches.append(current)
            current = []
        current.append(task)
    if current:
        batches.append(current)
    for i, batch in enumerate(batches, 1):
        name = f'batch-{i:03d}'
        folder = root / name
        folder.mkdir()
        save(folder / 'state.json', dict(id=name, status='pending', tasks=batch, attempts=0))
    return config


def folder_for(root, batch):
    if not re.fullmatch(r'batch-[0-9]{3,}', batch or '') or not (root / batch / 'state.json').is_file():
        raise ValueError('unknown batch')
    return root / batch


def snapshot(config, folder, state):
    work = folder / 'worktree'
    base = state.get('base', config['base'])
    if git(work, 'rev-parse', 'HEAD').decode().strip() != base:
        raise ValueError('author must not commit; worktree HEAD changed')
    allowed = {p for t in state['tasks'] for p in t['paths']}
    changed = set(git(work, 'diff', '--name-only', '-z', base).decode().strip('\0').split('\0'))
    changed.update(git(work, 'diff', '--cached', '--name-only', '-z', base).decode().strip('\0').split('\0'))
    changed.update(git(work, 'ls-files', '--others', '--exclude-standard', '-z').decode().strip('\0').split('\0'))
    changed.discard('')
    if changed - allowed:
        raise ValueError('out-of-scope paths: ' + ', '.join(sorted(changed - allowed)))
    for name in changed:
        path = work / name
        if path.is_symlink() or work.resolve() not in path.resolve().parents:
            raise ValueError('symlink/path escape: ' + name)
    if changed:
        git(work, 'add', '--all', '--', *sorted(changed))
    patch = git(work, 'diff', '--cached', '--binary', base)
    return patch, hashlib.sha256(patch).hexdigest()


def execute(config, folder, feedback=None):
    state = read(folder / 'state.json')
    work = folder / 'worktree'
    if feedback is None:
        if state['status'] != 'pending':
            raise ValueError('only pending batches may start')
    elif state['status'] not in ('changes-requested', 'failed', 'awaiting-review', 'approved'):
        raise ValueError('batch is not revisable')
    state.pop('approved_digest', None)
    state['attempts'] += 1
    state['status'] = 'running'
    save(folder / 'state.json', state)
    try:
        if feedback is None:
            if not clean(config['repo']):
                raise ValueError('dispatch target must be clean')
            state['base'] = git(config['repo'], 'rev-parse', 'HEAD').decode().strip()
        if feedback is None or not work.exists():
            git(config['repo'], 'worktree', 'add', '--detach', str(work), state.get('base', config['base']))
        command = [config['codex'], 'exec']
        if feedback is not None and state.get('session_id'):
            session = str(uuid.UUID(state['session_id']))
            command += ['resume', session]
        command += ['--json', '-m', config['model'], '-c',
                    'model_reasoning_effort=' + json.dumps(config['reasoning']),
                    '-c', 'sandbox_mode="workspace-write"', '-c', 'approval_policy="never"', '-']
        prompt = ('Act as an editorial author. Follow repository instructions, docs/content-review-checklist.md '
                  'and .agents/skills/editorial-author/SKILL.md if available. Do not take the coordinator role. '
                  'Implement every acceptance criterion below; run applicable checks. '
                  'Edit only the exact task paths. Do not commit, integrate, approve, spawn agents, or edit coordinator state. '
                  'Leave changes for primary review and report task IDs, evidence and remaining limitations.\n' +
                  json.dumps(state['tasks'], ensure_ascii=False))
        if feedback is not None:
            prompt += '\nPrimary reviewer requests these corrections:\n' + feedback
        log = folder / f'attempt-{state["attempts"]:03d}.jsonl'
        timed_out = False
        if not state.get('go_cache'):
            state['go_cache'] = tempfile.mkdtemp(prefix='codinho-editorial-go-')
        save(folder / 'state.json', state)
        env = dict(os.environ, GOCACHE=state['go_cache'], GOFLAGS='-mod=readonly')
        with log.open('wb') as output, log.with_suffix('.stderr').open('wb') as errors:
            proc = subprocess.Popen(command, cwd=work, stdin=subprocess.PIPE, stdout=output,
                                    stderr=errors, start_new_session=True, env=env)
            state['author_pid'] = proc.pid
            save(folder / 'state.json', state)
            try:
                proc.communicate(prompt.encode(), timeout=config['timeout_seconds'])
            except subprocess.TimeoutExpired:
                os.killpg(proc.pid, signal.SIGTERM)
                try:
                    proc.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    os.killpg(proc.pid, signal.SIGKILL)
                    proc.wait()
                timed_out = True
        for line in log.read_text().splitlines():
            try:
                event = json.loads(line)
                if event.get('type') == 'thread.started':
                    state['session_id'] = str(uuid.UUID(event['thread_id']))
            except (ValueError, KeyError):
                continue
        if timed_out:
            raise ValueError('author timed out; work and explicit session preserved')
        if proc.returncode:
            raise ValueError(f'author exited {proc.returncode}; see {log}')
        if not state.get('session_id'):
            raise ValueError('author did not emit a valid session ID')
        patch, digest = snapshot(config, folder, state)
        (folder / 'proposed.patch').write_bytes(patch)
        state.update(status='awaiting-review', digest=digest)
        state.pop('error', None)
    except (ValueError, OSError, subprocess.SubprocessError) as error:
        state.update(status='failed', error=str(error))
    save(folder / 'state.json', state)
    return state


def recover(config, folder):
    state = read(folder / 'state.json')
    if state['status'] != 'running':
        raise ValueError('recover requires an interrupted running batch')
    pid = state.get('author_pid')
    if not pid:
        raise ValueError('missing author PID; inspect legacy execution manually')
    try:
        os.kill(pid, 0)
    except ProcessLookupError:
        pass
    else:
        raise ValueError('author process still exists; recovery refused')
    log = folder / f'attempt-{state["attempts"]:03d}.jsonl'
    for line in log.read_text().splitlines():
        try:
            event = json.loads(line)
            if event.get('type') == 'thread.started':
                state['session_id'] = str(uuid.UUID(event['thread_id']))
        except (ValueError, KeyError):
            continue
    state.update(status='failed', error='interrupted execution recovered; revise required')
    save(folder / 'state.json', state)
    return state


def review(config, folder, decision, feedback):
    if decision not in ('approve', 'request-changes') or not feedback.strip():
        raise ValueError('review requires a valid decision and nonempty evidence')
    state = read(folder / 'state.json')
    if state['status'] != 'awaiting-review':
        raise ValueError('only delivered batches may be reviewed')
    patch, digest = snapshot(config, folder, state)
    if digest != state['digest']:
        raise ValueError('diff changed since delivery; request revision first')
    state.setdefault('reviews', []).append(dict(decision=decision, feedback=feedback, digest=digest))
    state['status'] = 'approved' if decision == 'approve' else 'changes-requested'
    if decision == 'approve':
        state['approved_digest'] = digest
    save(folder / 'state.json', state)
    return state


def integrate(config, folder):
    state = read(folder / 'state.json')
    if state['status'] != 'approved':
        raise ValueError('primary approval is required')
    patch, digest = snapshot(config, folder, state)
    if digest != state.get('approved_digest'):
        raise ValueError('stale approval: diff changed')
    repo = config['repo']
    if not clean(repo):
        raise ValueError('integration target must be clean')
    if not patch:
        base = state.get('base', config['base'])
        if git(repo, 'rev-parse', 'HEAD').decode().strip() != base:
            raise ValueError('no-change audit is stale; target HEAD changed')
        state.update(status='integrated', commit=base, no_changes=True)
        save(folder / 'state.json', state)
        return state
    patch_path = folder / 'approved.patch'
    patch_path.write_bytes(patch)
    state['integration_base'] = git(repo, 'rev-parse', 'HEAD').decode().strip()
    save(folder / 'state.json', state)
    git(repo, 'apply', '--check', '--index', str(patch_path))
    git(repo, 'apply', '--index', str(patch_path))
    state['integration_digest'] = hashlib.sha256(
        git(repo, 'diff', '--cached', '--binary', state['integration_base'])).hexdigest()
    save(folder / 'state.json', state)
    try:
        git(repo, 'commit', '-m', 'feat(content): editorial ' + state['id'],
            '-m', 'POSE-Spec: ' + config['queue']['spec'])
    except subprocess.SubprocessError:
        state['status'] = 'integration-failed'
        save(folder / 'state.json', state)
        raise ValueError('commit failed; staged changes preserved, complete integration manually')
    state.update(status='integrated', commit=git(repo, 'rev-parse', 'HEAD').decode().strip())
    save(folder / 'state.json', state)
    return state


def reconcile(config, folder):
    state = read(folder / 'state.json')
    if state['status'] != 'integration-failed' or not clean(config['repo']):
        raise ValueError('reconcile requires integration-failed and a clean target')
    repo = config['repo']
    head = git(repo, 'rev-parse', 'HEAD').decode().strip()
    parents = git(repo, 'rev-list', '--parents', '-n', '1', head).decode().split()[1:]
    if parents != [state.get('integration_base')]:
        raise ValueError('manual commit must directly follow integration base')
    patch = git(repo, 'diff', '--binary', state['integration_base'], head)
    if hashlib.sha256(patch).hexdigest() != state.get('integration_digest', state['approved_digest']):
        raise ValueError('manual commit differs from approved patch')
    message = git(repo, 'log', '-1', '--format=%B').decode().splitlines()
    if 'POSE-Spec: ' + config['queue']['spec'] not in message:
        raise ValueError('manual commit requires the matching POSE-Spec trailer')
    state.update(status='integrated', commit=head)
    save(folder / 'state.json', state)
    return state


def ready_batches(config, folders, limit):
    states = [(p.parent, read(p)) for p in folders]
    completed = set(config.get('completed_task_ids', []))
    completed.update(t['id'] for _, s in states if s['status'] == 'integrated' for t in s['tasks'])
    selected = []
    occupied = {p for _, s in states if s['status'] not in ('pending', 'integrated')
                for t in s['tasks'] for p in t['paths']}
    for folder, state in states:
        paths = {p for t in state['tasks'] for p in t['paths']}
        deps = {d for t in state['tasks'] for d in t.get('depends_on', [])}
        if state['status'] == 'pending' and deps <= completed and not paths & occupied:
            if len(selected) < limit:
                selected.append(folder)
        if state['status'] != 'integrated':
            occupied.update(paths)
    return selected


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest='command', required=True)
    for name in ('init', 'status', 'run', 'revise', 'review', 'integrate', 'recover', 'reconcile'):
        p = commands.add_parser(name)
        p.add_argument('--run-dir', required=True)
        if name == 'init':
            p.add_argument('--queue', required=True)
            p.add_argument('--model', default='gpt-5.6-luna')
            p.add_argument('--reasoning', choices=('low', 'medium', 'high', 'xhigh', 'max'), default='high')
            p.add_argument('--parallelism', type=int, default=1)
            p.add_argument('--batch-size', type=int, default=5)
            p.add_argument('--timeout-seconds', type=int, default=1800)
            p.add_argument('--codex', default='codex')
            p.add_argument('--completed-task', action='append', default=[])
            p.add_argument('--task', action='append', default=[])
        if name in ('revise', 'review', 'integrate', 'recover', 'reconcile'):
            p.add_argument('--batch', required=True)
        if name == 'run':
            p.add_argument('--limit-batches', type=int)
        if name in ('revise', 'review'):
            p.add_argument('--feedback', required=True)
        if name == 'review':
            p.add_argument('--decision', choices=('approve', 'request-changes'), required=True)
    args = parser.parse_args(argv)
    try:
        if args.command == 'init':
            result = init(args)
        else:
            root = Path(args.run_dir).resolve()
            config = read(root / 'config.json')
            if args.command == 'status':
                print(json.dumps([read(p) for p in sorted(root.glob('batch-*/state.json'))], ensure_ascii=False, indent=2))
                return 0
            if args.command in ('revise', 'recover'):
                folder = folder_for(root, args.batch)
                with (folder / 'author.lock').open('a') as lock:
                    fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
                    if args.command == 'revise':
                        result = execute(config, folder, args.feedback)
                    else:
                        result = recover(config, folder)
                print(json.dumps(result, ensure_ascii=False, indent=2))
                return int(result.get('status') == 'failed')
            with (root / 'coordinator.lock').open('a') as lock:
                fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
                folders = sorted(root.glob('batch-*/state.json'))
                if args.command == 'status':
                    result = [read(p) for p in folders]
                elif args.command == 'run':
                    limit = config['parallelism']
                    if args.limit_batches is not None:
                        if args.limit_batches < 1:
                            raise ValueError('limit-batches must be positive')
                        limit = min(limit, args.limit_batches)
                    pending = ready_batches(config, folders, limit)
                    with concurrent.futures.ThreadPoolExecutor(max_workers=config['parallelism']) as pool:
                        result = list(pool.map(lambda p: execute(config, p), pending))
                else:
                    folder = folder_for(root, args.batch)
                    if args.command == 'review':
                        result = review(config, folder, args.decision, args.feedback)
                    elif args.command == 'recover':
                        result = recover(config, folder)
                    elif args.command == 'reconcile':
                        result = reconcile(config, folder)
                    else:
                        result = integrate(config, folder)
        print(json.dumps(result, ensure_ascii=False, indent=2))
        states = result if isinstance(result, list) else [result]
        return int(any(s.get('status') == 'failed' for s in states))
    except (ValueError, OSError, subprocess.SubprocessError) as error:
        print(str(error), file=sys.stderr)
        return 1


if __name__ == '__main__':
    sys.exit(main())
