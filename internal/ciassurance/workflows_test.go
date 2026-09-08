package ciassurance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Parse the actual workflows so pinning applies to every job and step.
func TestWorkflowTrustBoundaries(t *testing.T) {
	for _, name := range []string{"ci.yml", "release.yml"} {
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("../../.github/workflows", name))
			if err != nil {
				t.Fatal(err)
			}
			var workflow struct {
				Permissions map[string]string `yaml:"permissions"`
				Jobs        map[string]struct {
					Steps []struct {
						Uses string `yaml:"uses"`
						Run  string `yaml:"run"`
					} `yaml:"steps"`
				} `yaml:"jobs"`
			}
			if err := yaml.Unmarshal(raw, &workflow); err != nil {
				t.Fatal(err)
			}
			if workflow.Permissions["contents"] != "read" {
				t.Fatal("workflow must only read repository content")
			}
			pins := regexp.MustCompile(`^[\w/-]+@[a-f0-9]{40}$`)
			if len(workflow.Jobs) == 0 {
				t.Fatal("no jobs")
			}
			for job, body := range workflow.Jobs {
				for _, step := range body.Steps {
					if step.Uses != "" && !pins.MatchString(step.Uses) {
						t.Errorf("%s has unpinned action %s", job, step.Uses)
					}
					if strings.Contains(step.Run, "${{") {
						t.Errorf("%s interpolates expressions into shell", job)
					}
				}
			}
		})
	}
}

func TestReleaseRejectsUnsafeLabelsBeforeBuild(t *testing.T) {
	for _, label := range []string{"", "1.2.3", "v01.2.3", "v1.2.3;touch marker", "v1.2.3$(touch marker)", "v1.2.3\nwhoami", "v1.2.3 -X other=value", strings.Repeat("v", 81)} {
		t.Run(label, func(t *testing.T) {
			cmd := exec.Command("bash", "../../scripts/ci/build-release.sh", label)
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("accepted invalid label")
			}
			if !strings.Contains(string(out), "Invalid build version") {
				t.Fatalf("did not reject before build: %s (%v)", out, err)
			}
		})
	}
}

// Large version output must not trigger SIGPIPE under set -o pipefail.
func TestValidationReachesGatesAfterVersionCheck(t *testing.T) {
	bin := t.TempDir()
	scripts := map[string]string{
		"pose":        "#!/bin/bash\ncase \"$1\" in\nversion) printf 'pose 1.7.12\\n'; printf '%100000s\\n' details ;;\nskills-check) echo reached-skills-gate; exit 23 ;;\n*) exit 99 ;;\nesac\n",
		"go":          "#!/bin/bash\nprintf 'mod golang.org/x/vuln v1.6.0 h1:FeMO9Rm/HwyduOztbvKcOw+zvDEPr4I4aQNSfevFcKY=\\n'\n",
		"govulncheck": "#!/bin/bash\nexit 99\n",
	}
	for name, body := range scripts {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(body), 0755); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("bash", "../../scripts/ci/validate.sh")
	cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "reached-skills-gate") {
		t.Fatalf("gate not reached: %s (%v)", out, err)
	}
	if status, ok := err.(*exec.ExitError); !ok || status.ExitCode() != 23 {
		t.Fatalf("gate failure was not propagated: %v", err)
	}
}

func TestHistoricalRenameMatchesGit(t *testing.T) {
	raw, err := os.ReadFile("../../.pose/contracts/historical-renames.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		SchemaVersion int `json:"schema_version"`
		Renames       []struct {
			Commit    string `json:"commit"`
			OldPath   string `json:"old_path"`
			NewPath   string `json:"new_path"`
			GitStatus string `json:"git_status"`
		} `json:"renames"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != 1 || len(manifest.Renames) == 0 {
		t.Fatal("missing historical rename evidence")
	}
	for _, rename := range manifest.Renames {
		if !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(rename.Commit) {
			t.Fatal("rename requires immutable commit")
		}
		cmd := exec.Command("git", "show", "--format=", "--name-status", "--find-renames", rename.Commit, "--", rename.OldPath, rename.NewPath)
		cmd.Dir = "../.."
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		want := rename.GitStatus + "\t" + rename.OldPath + "\t" + rename.NewPath
		if strings.TrimSpace(string(out)) != want {
			t.Fatalf("rename differs from Git: %s", out)
		}
	}
}

func TestGovernancePoliciesRequireRealEntrypoints(t *testing.T) {
	raw, err := os.ReadFile("../../.pose/policy/delivery.json")
	if err != nil {
		t.Fatal(err)
	}
	var policy struct {
		Enabled bool `json:"enabled"`
		Roots   []struct {
			Path       string `json:"path"`
			Entrypoint string `json:"entrypoint"`
		} `json:"roots"`
	}
	if err := json.Unmarshal(raw, &policy); err != nil {
		t.Fatal(err)
	}
	if !policy.Enabled || len(policy.Roots) == 0 {
		t.Fatal("delivery policy must govern real roots")
	}
	for _, root := range policy.Roots {
		if !filepath.IsLocal(root.Path) || !filepath.IsLocal(root.Entrypoint) {
			t.Fatal("delivery paths must remain local")
		}
		relative, err := filepath.Rel(root.Path, root.Entrypoint)
		if err != nil || !filepath.IsLocal(relative) {
			t.Fatal("entrypoint outside root")
		}
		if info, err := os.Stat(filepath.Join("../..", root.Entrypoint)); err != nil || info.IsDir() {
			t.Fatalf("missing entrypoint %s", root.Entrypoint)
		}
	}
}
