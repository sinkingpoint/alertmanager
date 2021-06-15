// Copyright 2019 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wechat

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-kit/kit/log"
	commoncfg "github.com/prometheus/common/config"
	"github.com/stretchr/testify/require"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify/test"
)

func TestWechatRedactedURLOnInitialAuthentication(t *testing.T) {
	ctx, u, fn := test.GetContextWithCancelingURL()
	defer fn()

	secret := "secret_key"
	notifier, err := New(
		&config.WechatConfig{
			APIURL:     &config.URL{URL: u},
			HTTPConfig: &commoncfg.HTTPClientConfig{},
			CorpID:     "corpid",
			APISecret:  config.Secret(secret),
		},
		test.CreateTmpl(t),
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	test.AssertNotifyLeaksNoSecret(t, ctx, notifier, secret)
}

func TestWechatRedactedURLOnNotify(t *testing.T) {
	secret, token := "secret", "token"
	ctx, u, fn := test.GetContextWithCancelingURL(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"access_token":"%s"}`, token)
	})
	defer fn()

	notifier, err := New(
		&config.WechatConfig{
			APIURL:     &config.URL{URL: u},
			HTTPConfig: &commoncfg.HTTPClientConfig{},
			CorpID:     "corpid",
			APISecret:  config.Secret(secret),
		},
		test.CreateTmpl(t),
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	test.AssertNotifyLeaksNoSecret(t, ctx, notifier, secret, token)
}

func TestWechatAPISecretFromFile(t *testing.T) {
	secret, token := "secret", "token"
	f, err := ioutil.TempFile("", "victorops_test")
	require.NoError(t, err, "creating temp file failed")
	_, err = f.WriteString(secret)
	require.NoError(t, err, "writing to temp file failed")

	ctx, u, fn := test.GetContextWithCancelingURL(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"access_token":"%s"}`, token)
	})
	defer fn()

	notifier, err := New(
		&config.WechatConfig{
			APIURL:        &config.URL{URL: u},
			HTTPConfig:    &commoncfg.HTTPClientConfig{},
			CorpID:        "corpid",
			APISecretFile: f.Name(),
		},
		test.CreateTmpl(t),
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	test.AssertNotifyLeaksNoSecret(t, ctx, notifier, secret, token)
}

func TestWechatMessageTypeSelector(t *testing.T) {
	secret, token := "secret", "token"
	ctx, u, fn := test.GetContextWithCancelingURL(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"access_token":"%s"}`, token)
	})
	defer fn()

	notifier, err := New(
		&config.WechatConfig{
			APIURL:      &config.URL{URL: u},
			HTTPConfig:  &commoncfg.HTTPClientConfig{},
			CorpID:      "corpid",
			APISecret:   config.Secret(secret),
			MessageType: "markdown",
		},
		test.CreateTmpl(t),
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	test.AssertNotifyLeaksNoSecret(t, ctx, notifier, secret, token)
}
