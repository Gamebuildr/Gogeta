package main

import "testing"

func TestGoRepoPathAlwaysUnique(t *testing.T) {
	testRepoPath, testRelativePath := GetRepoPath("testProject")
	copyRepoPathCopy, copyRelativePath := GetRepoPath("testProject")
	if testRepoPath == copyRepoPathCopy {
		t.Fatalf("Expected repoPath to not be equal")
	}
	if testRelativePath == copyRelativePath {
		t.Fatalf("Expected relativePath to not be equal")
	}
}
