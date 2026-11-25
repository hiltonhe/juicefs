/*
 * JuiceFS, Copyright 2026 Juicedata, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package object

import "testing"

func TestApplyFormatCredentialOverride(t *testing.T) {
	withCreds := func(storage string) VolumeConfig {
		return VolumeConfig{
			Storage:      storage,
			AccessKey:    "ak",
			SecretKey:    "sk",
			SessionToken: "token",
		}
	}
	assertEmpty := func(t *testing.T, c VolumeConfig) {
		t.Helper()
		if c.AccessKey != "" || c.SecretKey != "" || c.SessionToken != "" {
			t.Fatalf("expected credentials to be cleared, got %+v", c)
		}
	}

	t.Run("disabled", func(t *testing.T) {
		t.Setenv(ignoreFormatCredentialsEnv, "0")
		got := applyFormatCredentialOverride(withCreds("s3"))
		if got.AccessKey != "ak" || got.SecretKey != "sk" || got.SessionToken != "token" {
			t.Fatalf("credentials should be kept, got %+v", got)
		}
	})

	t.Run("unset", func(t *testing.T) {
		t.Setenv(ignoreFormatCredentialsEnv, "")
		got := applyFormatCredentialOverride(withCreds("oss"))
		if got.AccessKey != "ak" {
			t.Fatalf("credentials should be kept, got %+v", got)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv(ignoreFormatCredentialsEnv, "not-a-bool")
		got := applyFormatCredentialOverride(withCreds("s3"))
		if got.AccessKey != "ak" {
			t.Fatalf("credentials should be kept, got %+v", got)
		}
	})

	t.Run("s3", func(t *testing.T) {
		t.Setenv(ignoreFormatCredentialsEnv, "true")
		assertEmpty(t, applyFormatCredentialOverride(withCreds("s3")))
	})

	t.Run("oss-uppercase", func(t *testing.T) {
		t.Setenv(ignoreFormatCredentialsEnv, "1")
		assertEmpty(t, applyFormatCredentialOverride(withCreds("OSS")))
	})

	t.Run("unsupported-storage", func(t *testing.T) {
		t.Setenv(ignoreFormatCredentialsEnv, "true")
		got := applyFormatCredentialOverride(withCreds("file"))
		if got.AccessKey != "ak" || got.SecretKey != "sk" {
			t.Fatalf("unsupported storage should keep static credentials, got %+v", got)
		}
	})
}
