package config

import (
	"testing"
	"os"
	"runtime"
)

const (
	testIP   = "1.2.3.4"
	testPort = "7532"
)

func TestGetEnvOrDefault(t *testing.T) {
	type args struct {
		envKey       string
		defaultValue string
	}
	const (
		testValDefault = "ThisIsATest"
		testValEnv     = "MyValue"
		testKeyEnv     = "MyKey"
	)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "unknown_env_key",
			args: args{
				envKey:       "UNKNOWN_ENV_KEY",
				defaultValue: testValDefault,
			},
			want: testValDefault,
		},
		{
			name: "existing_env_key",
			args: args{
				envKey:       testKeyEnv,
				defaultValue: testValEnv,
			},
			want: testValEnv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(testKeyEnv, testValEnv)
			if got := getEnvOrDefault(tt.args.envKey, tt.args.defaultValue); got != tt.want {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAddr(t *testing.T) {
	type envVars map[string]string

	tests := []struct {
		name string
		env  envVars
		want string
	}{
		{
			name: "default_ip_and_port",
			env:  envVars{},
			want: defaultServerIP + ":" + defaultServerPort,
		},
		{
			name: "default_port",
			env: envVars{
				serverIPKey: testIP,
			},
			want: testIP + ":" + defaultServerPort,
		},
		{
			name: "default_ip",
			env: envVars{
				serverPortKey: testPort,
			},
			want: defaultServerIP + ":" + testPort,
		},
		{
			name: "ok",
			env: envVars{
				serverIPKey:   testIP,
				serverPortKey: testPort,
			},
			want: testIP + ":" + testPort,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			extractConfFromEnvVars() // force re-extraction
			if got := ServerAddr(); got != tt.want {
				t.Errorf("ServerAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListenAddr(t *testing.T) {
	type envVars map[string]string

	tests := []struct {
		name string
		env  envVars
		want string
	}{
		{
			name: "default_ip_and_port",
			env:  envVars{},
			want: defaultListenIP + ":" + defaultListenPort,
		},
		{
			name: "default_port",
			env: envVars{
				listenIPKey: testIP,
			},
			want: testIP + ":" + defaultListenPort,
		},
		{
			name: "default_ip",
			env: envVars{
				listenPortKey: testPort,
			},
			want: defaultListenIP + ":" + testPort,
		},
		{
			name: "ok",
			env: envVars{
				listenIPKey:   testIP,
				listenPortKey: testPort,
			},
			want: testIP + ":" + testPort,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			extractConfFromEnvVars() // force re-extraction
			if got := ListenAddr(); got != tt.want {
				t.Errorf("ListenAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestWorkerCount(t *testing.T) {
	want := runtime.NumCPU()
	if got := WorkerCount(); got != want {
		t.Errorf("WorkerCount() = %v, want %v", got, want)
	}
}
