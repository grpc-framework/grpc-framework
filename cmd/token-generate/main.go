/*
Copyright © 2020 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/grpc-framework/grpc-framework/v2/pkg/auth"
	"github.com/sirupsen/logrus"
)

type tokenInfo struct {
	issuer  string
	subject string
	name    string
	email   string
	roles   string
	groups  string
}

type tokenGenOptions struct {
	sharedSecret string
	rsaPem       string
	ecdsaPem     string
	duration     string
	output       string
	token        tokenInfo
}

func main() {
	tokenGenArgs := &tokenGenOptions{
		token: tokenInfo{},
	}

	// Define flags
	flag.StringVar(&tokenGenArgs.sharedSecret, "shared-secret", "", "Shared secret to sign token")
	flag.StringVar(&tokenGenArgs.rsaPem, "rsa-private-keyfile", "", "RSA Private file to sign token")
	flag.StringVar(&tokenGenArgs.ecdsaPem, "ecdsa-private-keyfile", "", "ECDSA Private file to sign token")
	flag.StringVar(&tokenGenArgs.duration, "token-duration", "1d", "Duration of time where the token will be valid. "+
		"Postfix the duration by using s for seconds, m for minutes, h for hours, d for days, and y for years.")
	flag.StringVar(&tokenGenArgs.token.issuer, "token-issuer", "portworx.com",
		"Issuer name of token. Do not use https:// in the issuer since it could indicate that this is an OpenID Connect issuer.")
	flag.StringVar(&tokenGenArgs.token.name, "token-name", "", "Account name")
	flag.StringVar(&tokenGenArgs.token.subject, "token-subject", "", "Unique ID of this account")
	flag.StringVar(&tokenGenArgs.token.email, "token-email", "", "Email address of the account")
	flag.StringVar(&tokenGenArgs.token.roles, "token-roles", "", "Comma separated list of roles applied to this token")
	flag.StringVar(&tokenGenArgs.token.groups, "token-groups", "", "Comma separated list of groups which the token will be part of")

	// Parse flags
	flag.Parse()

	// Validate required flags
	if len(tokenGenArgs.token.name) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Must supply an account name")
		os.Exit(1)
	}
	if len(tokenGenArgs.token.email) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Must supply an email address")
		os.Exit(1)
	}
	if len(tokenGenArgs.token.subject) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Must supply a unique identifier as the subject")
		os.Exit(1)
	}
	if len(tokenGenArgs.token.roles) == 0 {
		logrus.Warning("Warning: No role provided")
	}
	if len(tokenGenArgs.token.groups) == 0 {
		logrus.Warning("Warning: No groups provided")
	}

	claims := &auth.Claims{
		Name:    tokenGenArgs.token.name,
		Email:   tokenGenArgs.token.email,
		Subject: tokenGenArgs.token.subject,
		Roles:   strings.Split(tokenGenArgs.token.roles, ","),
		Groups:  strings.Split(tokenGenArgs.token.groups, ","),
		Issuer:  tokenGenArgs.token.issuer,
	}

	// Get duration
	expDuration, err := auth.ParseToDuration(tokenGenArgs.duration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to parse duration\n")
		os.Exit(1)
	}

	options := &auth.Options{
		Expiration: time.Now().Add(expDuration).Unix(),
	}

	// Get signature
	var signature *auth.Signature
	if len(tokenGenArgs.sharedSecret) != 0 {
		signature, err = auth.NewSignatureSharedSecret(tokenGenArgs.sharedSecret)
	} else if len(tokenGenArgs.rsaPem) != 0 {
		signature, err = auth.NewSignatureRSAFromFile(tokenGenArgs.rsaPem)
	} else if len(tokenGenArgs.ecdsaPem) != 0 {
		signature, err = auth.NewSignatureECDSAFromFile(tokenGenArgs.ecdsaPem)
	} else {
		fmt.Fprintln(os.Stderr, "Error: Must provide a secret key to sign token")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate signature: %v\n", err)
		os.Exit(1)
	}

	// Generate token
	token, err := auth.Token(claims, signature, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create token: %v\n", err)
		os.Exit(1)
	}

	// Print token
	fmt.Println(token)
}
