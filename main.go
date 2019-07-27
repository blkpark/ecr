package main

import (
  "os"
  "fmt"
  "time"
  "errors"
  "strings"
  "encoding/base64"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ecr"
)

const (
  RETRY = 1
)

type parameter struct {
  key string
  secret string
  region string
  registry string
  retry int
}

type auth struct {
	username string
	password string
	proxy string
	expiration time.Time
}

func env(key, val string) string {
	v := os.Getenv(key)
	if v == "" {
		return val
	} else {
		return v
	}
}

func getParameter() (*parameter, error) {
  key := env("KEY", "")
  secret := env("SECRET", "")
  region := env("REGION", "")
  registry := env("REGISTRY", "")

  var err error

  for {
    if key == "" {
      err = errors.New("Invalid parameter for KEY")
      break
    }

    if secret == "" {
      err = errors.New("Invalid parameter for SECRET")
      break
    }

    if region == "" {
      err = errors.New("Invalid parameter for REGION")
      break
    }

    if registry == "" {
      err = errors.New("Invalid parameter for REGISTRY")
      break
    }

    break // success
  }

  if err != nil {
    return nil, err
  }

  p := &parameter{
    key: key,
    secret: secret,
    region: region,
    registry: registry,
    retry: RETRY,
  }

  return p, nil
}

func getAuths(p *parameter) ([]*auth, error) {
  s, err := session.NewSession(&aws.Config{
    Region: aws.String(p.region),
    Credentials: credentials.NewStaticCredentials(p.key, p.secret, ""),
  })

  if err != nil {
    return nil, err
  }

  input := &ecr.GetAuthorizationTokenInput{
    RegistryIds: []*string{aws.String(p.registry)},
  }

	e := ecr.New(s, aws.NewConfig().WithMaxRetries(p.retry).WithRegion(p.region))

  res, err := e.GetAuthorizationToken(input)
  if err != nil {
    return nil, errors.New("No token")
  }

  auths := []*auth{}

	for _, v := range res.AuthorizationData {
		token, err := base64.StdEncoding.DecodeString(*v.AuthorizationToken)
    if err != nil {
      break
    }

		splited := strings.SplitN(string(token), ":", 2)
    if len(splited) != 2 {
      err = errors.New("Found wrong format")
      break
    }

    a := &auth{
      username: splited[0],
      password: splited[1],
			proxy: *(v.ProxyEndpoint),
			expiration: *(v.ExpiresAt),
    }

    auths = append(auths, a)
	}

  if err != nil {
    return nil, err
  }

  if len(auths) == 0 {
    return nil, errors.New("No auths")
  }

  return auths, nil
}

func main() {
  p, err := getParameter()
  if err != nil {
    return
  }

  auths, err := getAuths(p)
  if err != nil {
    return
  }

  a := auths[0]
  cmd := fmt.Sprintf("docker login -u %s -p %s %s", a.username, a.password, a.proxy)
  fmt.Println(cmd)
}
