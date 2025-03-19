package main

type sessionKey string

const postLoginRedirectURLSessionKey = sessionKey("postLoginRedirectURL")
const authenticatedUserIDSessionKey = sessionKey("authenticatedUserID")
