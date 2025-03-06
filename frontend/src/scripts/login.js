"use strict";

import auth from "./common/authorize.js"
import apiclient from "./common/apis.js"
import words from "./common/i18n.js"

(()=>{
  // salt for password for each login session
  var se_token="";
  // initialize API client
  var api = new apiclient(window.baseurl?window.baseurl:"");
  // flag for allowing commit
  var commitmode = "none";
  // flag for remember me
  var remember = false;

  // show error message
  var show_error=(msg)=>{
    let err = document.getElementById("error-message");
    err.innerText = msg;
    err.style.display = "block";
  }

  var domref = {};

  var set_login_available = (mode)=>{
    switch(mode){
      case "ok":
        commitmode = "login";
        domref.error.style.display = "none";
        domref.login.disabled = false;
        domref.login.innerText = words.login;
        domref.user.disabled = false;
        domref.user.placeholder = words.username;
        domref.passwd.disabled = false
        domref.passwd.placeholder = words.password;
        domref.remember.disabled = false;
        break;
      case "failed":
        commitmode = "reflesh";
        domref.login.disabled = false;
        domref.login.innerText = words.reflesh;
        domref.user.disabled = true;
        domref.passwd.disabled = true;
        domref.passwd.value = "";
        domref.remember.disabled = true;
        break;
      case "waiting":
        commitmode = "none";
        domref.error.style.display = "none";
        domref.login.disabled = true;
        domref.login.innerText = words.waiting;
        domref.user.disabled = true;
        domref.passwd.disabled = true;
        domref.remember.disabled = true;
        break;
    }
  }

  // renew se_token from server
  var renew_setoken=()=>{
    set_login_available("waiting");
    api.get_login_token().then((res)=>{
      if (res.success()) {
        se_token = res.data;
        set_login_available("ok");
      }else{
        show_error("failed to get token: "+res.error);
        set_login_available("failed");
      }
    }).catch((e)=>{
      show_error("failed to get token: "+e);
      set_login_available("failed");
    });
  }

  // commit login
  var commit= async (user, passwd)=>{
    remember = domref.remember.checked;
    const use_se = await auth.get_user_se(user);
    const pwd_se = await auth.get_pwd_se(passwd, use_se, se_token);
    //TODO commit
    console.log({
      user: user,
      passwd: passwd,
      use_se: use_se,
      pwd_se: pwd_se,
    });
    
  }

  document.addEventListener("DOMContentLoaded", ()=>{
    let i18ntags = document.querySelectorAll("[i18nword]");
    i18ntags.forEach((tag)=>{
      let word = tag.getAttribute("i18nword");
      tag.innerText = words[word];
    })
    domref["user"] = document.getElementById("user");
    domref["passwd"] = document.getElementById("passwd");
    domref["login"] = document.getElementById("login");
    domref["error"] = document.getElementById("error-message");
    domref["remember"] = document.getElementById("remember");
    domref["login"].addEventListener("click", ()=>{
      if (commitmode == "reflesh") {
        renew_setoken();
        return;
      }else if (commitmode != "login") {
        return;
      }
      let user = domref["user"].value;
      let passwd = domref["passwd"].value;
      if (!user || !passwd) {
        show_error("username or password is empty");
        return;
      }
      // TODO: lock the form
      commit(user, passwd).then((token)=>{
        window.localStorage.setItem("token", token);
        window.location.href = "subpages/ctrplane.html";
        //TODO: show result
      }).catch((e)=>{
        //TODO: show result
      }).finally(()=>{
        //TODO: unlock the form
      })
    })

    // initialize
    renew_setoken();
  })
})()
