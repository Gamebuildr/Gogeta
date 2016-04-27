package main

import "testing"

func TestGitClone(t *testing.T) {
    service := gogetaService{};
    expected := "Git Clone"
    actual := service.GitClone("test2")
    if(actual != expected) {
        t.Error("Git Clone Test Failed")
    }
}
