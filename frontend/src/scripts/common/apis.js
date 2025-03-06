"use strict";

// list of API urls
const urlist = {
  "token": "/api/token",
  "login": "/api/login",
  "prof": "/api/perf",
  "op": "/api/op"
}

// Result represents the result of an API call
class Result {
  constructor(obj, err){
    if (err) {
      this.error = err;
    }else if (!obj) {
      this.error = "No data received";
    }else if (obj.status === "error") {
      if (obj.error) {
        this.error = obj.error;
      }else{
        this.error = "Unknown error";
      }
    }else{
      this.data = obj.body;
    }
  }

  // check if the result is successful
  success() {
    return !this.error;
  }
}

// common fetch function
async function fetch_result(url, token, body) {
  let method = 'GET';
  let headers = {};
  if (body) {
    method = 'POST';
    headers['Content-Type'] = 'application/json';
  }
  if (token) {
    headers['WSDC-Token'] = token;
  }
  return await fetch(url, {
    method: method,
    headers: headers,
    body: body
  }).then((res) => {
    if (!res.ok || res.status !== 200) {
      return Promise.reject(res);
    }
    return res.json();
  }).catch((err) => {
    return new Result(null, err);
  }).then((obj) => {
    if (obj instanceof Result) {
      return obj;
    }
    return new Result(obj);
  })
}

// Api represents the API client
class Api {
  constructor(baseurl){
    this.baseurl = baseurl;
  }

  // get the login token
  async get_login_token() {
    return await fetch_result(this.baseurl + urlist.token);
  }

  // do login request
  async login(authdata, token) {
    return await fetch_result(this.baseurl + urlist.login, null, JSON.stringify({
      "authdata": authdata,
      "token": token
    }));
  }

  // get the performance data
  async get_perf(token) {
    return await fetch_result(this.baseurl + urlist.prof, token);
  }

  // call operation command
  async do_op(token, cmd, ...params) {
    return await fetch_result(this.baseurl + urlist.op, token, JSON.stringify({
      "cmd": cmd,
      "params": params
    }));
  }
}

export default Api;
