package handlers

import (
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
	"os/exec"
	"os"
	"bytes"
	"strconv"
)

func CriuPreDump(w http.ResponseWriter, r *http.Request) {
	var req RequestPreDump
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

	args := []string {
		"pre-dump",
		"--tree", strconv.Itoa(req.Pid),
	}
	if len(req.Path) > 0 {
		args = append(args, "--images-dir", req.Path)
	}
	if req.TrackMem {
		args = append(args, "--track-mem")
	}

	os.MkdirAll(req.Path, os.ModePerm)
	cmd := exec.Command("criu", args...)
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
	if err := json.NewEncoder(w).Encode(jsonOk{"PreDump done."}); err != nil {
		panic(err)
	}
	return
}

func CriuDump(w http.ResponseWriter, r *http.Request) {
	var req RequestDump
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

	args := []string {
		"dump",
		"--leave-stopped",
		"--tree", strconv.Itoa(req.Pid),
	}
	if len(req.Path) > 0 {
		args = append(args, "--images-dir", req.Path)
	}
	if len(req.PrevPath) > 0 {
		args = append(args, "--prev-images-dir", req.PrevPath)
	}
	if req.TrackMem {
		args = append(args, "--track-mem")
	}
	if req.Lazy {
		args = append(args, "--lazy-pages")
		args = append(args, "--address", "0.0.0.0")
		args = append(args, "--port", strconv.Itoa(req.LazyPort))
		args = append(args, "--status-fd", "1")				// STDOUT
	}
	if req.ShellJob {
		args = append(args, "--shell-job")
	}

	os.MkdirAll(req.Path, os.ModePerm)
	cmd := exec.Command("criu", args...)
	var b bytes.Buffer
	cmd.Stderr = &b
	printCommand(cmd)

	var stdout io.ReadCloser
	if req.Lazy {
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
	} else {
		cmd.Stdout = &b
	}

	if err := cmd.Start(); err != nil {
		printError(err)
	}

	if req.Lazy {
		b, _ := ioutil.ReadAll(stdout)
		if len(b) == 1 && b[0] == 0 {
			// ok
		} else {
			panic("Invalid status-fd return")
		}
	} else {
		if err := cmd.Wait(); err != nil {
			printError(err)
		}
		printOutput(b.Bytes())
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(jsonOk{"Dump done."}); err != nil {
		panic(err)
	}
}

func CriuRestore(w http.ResponseWriter, r *http.Request) {
	var req RequestRestore
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

	args := []string {
		"restore",
	}
	if len(req.Path) > 0 {
		args = append(args, "--images-dir", req.Path)
	}
	if req.Lazy {
		args = append(args, "--lazy-pages")
	}
	if req.ShellJob {
		args = append(args, "--shell-job")
	}

	cmd := exec.Command("criu", args...)
	cmd.Dir = req.Path
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
	if err := json.NewEncoder(w).Encode(jsonOk{"Restore done."}); err != nil {
		panic(err)
	}
	return
}

func CriuLazyPages(w http.ResponseWriter, r *http.Request) {
	var req RequestLazyPages
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

	args := []string {
		"lazy-pages",
		"--page-server",
		"--address", req.Address,
		"--port", strconv.Itoa(req.Port),
		"--status-fd", "1",					// STDOUT
	}

	cmd := exec.Command("criu", args...)
	cmd.Dir = req.Path
	var b bytes.Buffer
	cmd.Stderr = &b
	printCommand(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		printError(err)
	}

	if err := cmd.Start(); err != nil {
		printError(err)
	}

	b2, _ := ioutil.ReadAll(stdout)
	if len(b2) == 1 && b2[0] == 0 {
		// ok
	} else {
		panic("Invalid status-fd return")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(jsonOk{"Lazy-page client started."}); err != nil {
		panic(err)
	}
}