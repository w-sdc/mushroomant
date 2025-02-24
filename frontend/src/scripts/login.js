"use strict";

(()=>{
  // salt for public
  var public_salt="palworldwsdc@2025_"
  // salt for password for each login session
  var se_token="";

  // renew se_token from server
  var renew_setoken=()=>{
    // TODO
    console.log("renew_setoken yet not implemented")
  }

  // hash string to sha256 hex
  var strhash = async(s)=>{
    let bins = (new TextEncoder()).encode(s);
    return crypto.subtle.digest("sha-256", bins)
      .then((hbin)=>
        [...new Uint8Array(hbin)]
        .map((x)=>x.toString(16).padStart(2, '0')).join(""));
  }

  // commit login
  var commit= async (user, passwd)=>{
    let use_se = await strhash(public_salt + user)
    let pwd_se = await strhash([public_salt,use_se,passwd].join(";"))
    let pwd_se_ret = await strhash(se_token + pwd_se)
    //TODO commit
    console.log({
      user: user,
      passwd: passwd,
      use_se: use_se,
      pwd_se: pwd_se,
      pwd_se_ret: pwd_se_ret
    })
  }

  // initialize
  renew_setoken();

  // add event listener
  document.getElementById("login").addEventListener("click", ()=>{
    let user = document.getElementById("user").value;
    let passwd = document.getElementById("passwd").value;
    // TODO: lock the form
    commit(user, passwd).then(()=>{
      //TODO: show result
    }).catch((e)=>{
      //TODO: show result
    }).finally(()=>{
      //TODO: unlock the form
    })
  })
})()
