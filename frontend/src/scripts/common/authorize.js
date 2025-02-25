"use strict";

// salt for public
var public_salt="cdswdlrowlap@5202"

class AuthDigest {
  constructor(){
    this.tenc_utf8 = new TextEncoder();
  }

  // hash string to sha256 hex
  async strhash(s){
    return crypto.subtle.digest("sha-256", this.tenc_utf8.encode(s))
      .then((hbin)=>
        [...new Uint8Array(hbin)]
        .map((d)=>d.toString(16).padStart(2, '0')).join(""))
  }

  // get summary for user
  async get_user_se(user){
    return await this.strhash([public_salt,user].join(";"))
  }

  // get summary for password
  async get_pwd_se(passwd, user_se, se_token){
    let pwd_se = await this.strhash([public_salt,passwd].join(";"))
    return await this.strhash([se_token,user_se,pwd_se].join(";"))
  }
}

export default new AuthDigest()
