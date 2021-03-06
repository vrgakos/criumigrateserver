package handlers

import (
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
	"os/exec"
	"bytes"
)

func BaseScp(w http.ResponseWriter, r *http.Request) {
	var req RequestBaseScp
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	cmd := exec.Command("scp", "-r", req.From, req.To)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	printCommand(cmd)

	if err := cmd.Start(); err != nil {
		printError(err)
	}
	if err := cmd.Wait(); err != nil {
		printError(err)
	}
	printOutput(b.Bytes())

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(jsonOk{"Scp done."}); err != nil {
		panic(err)
	}
	return
}

func BaseClean(w http.ResponseWriter, r *http.Request) {
	var req RequestBaseClean
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	cmd := exec.Command("rm", "-rf", req.Path)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	printCommand(cmd)

	if err := cmd.Start(); err != nil {
		printError(err)
	}
	if err := cmd.Wait(); err != nil {
		printError(err)
	}
	printOutput(b.Bytes())

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(jsonOk{"Clean done."}); err != nil {
		panic(err)
	}
}