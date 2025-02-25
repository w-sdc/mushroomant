"use strict";

var words = {
  "en": {
    "login": "Login",
    "username": "Username",
    "password": "Password",
    "waiting": "Waiting...",
    "reflesh": "Reflesh",
    "savelogin": "Save Login",
  },
  "zh": {
    "login": "登录",
    "username": "用户名",
    "password": "密码",
    "waiting": "请稍后...",
    "reflesh": "刷新",
    "savelogin": "保存登录",
  },
};

var langremap = {
  "en": "en",
  "en-us": "en",
  "en-gb": "en",
  "zh": "zh",
  "zh-cn": "zh",
  "zh-hans-cn": "zh",
};

var current_lang = (()=>{
  let lang = document.documentElement.lang.toLowerCase();
  if (langremap[lang]) {
    return langremap[lang];
  }else{
    return "en";
  }
})();

export default words[current_lang];
