export function removeRegionCode(phone) {
  let res = phone.replace("+", "");
  if (res.indexOf("86") === 0) {
    res = res.replace("86", "");
  }
  return res;
}

export function getMixinContext() {
  let ctx = {};
  if (
    window.webkit &&
    window.webkit.messageHandlers &&
    window.webkit.messageHandlers.MixinContext
  ) {
    ctx = JSON.parse(prompt("MixinContext.getContext()"));
    ctx.platform = ctx.platform || "iOS";
  } else if (
    window.MixinContext &&
    typeof window.MixinContext.getContext === "function"
  ) {
    ctx = JSON.parse(window.MixinContext.getContext());
    ctx.platform = ctx.platform || "Android";
  }
  return ctx;
}

export function isMixin() {
  let mixinContext = getMixinContext();
  return mixinContext.app_version;
}

export function getPlatformType() {
  return isMixin() ? 0 : 1;
}

export function numberFmt(value) {
  return parseFloat(value).toLocaleString();
}

Date.prototype.Format = function (fmt) {
  var o = {
    "M+": this.getMonth() + 1, //month
    "d+": this.getDate(), //day
    "h+": this.getHours(), //Hour
    "m+": this.getMinutes(), //Minute
    "s+": this.getSeconds(), //Second
    "q+": Math.floor((this.getMonth() + 3) / 3), //Quarterly
    S: this.getMilliseconds(), //millisecond
  };
  if (/(y+)/.test(fmt))
    fmt = fmt.replace(
      RegExp.$1,
      (this.getFullYear() + "").substr(4 - RegExp.$1.length)
    );
  for (var k in o)
    if (new RegExp("(" + k + ")").test(fmt))
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length)
      );
  return fmt;
};

export function setHtmlMeta(name, content) {
  var metaList = document.getElementsByTagName("meta");
  for (var i = 0; i < metaList.length; i++) {
    if (metaList[i].getAttribute("name") === name) {
      metaList[i].content = content;
    }
  }
}

export function uuid() {
  var s = [];
  var hexDigits = "0123456789abcdef";
  for (var i = 0; i < 36; i++) {
    s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
  }
  s[14] = "4"; // bits 12-15 of the time_hi_and_version field to 0010
  s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1); // bits 6-7 of the clock_seq_hi_and_reserved to 01
  s[8] = s[13] = s[18] = s[23] = "-";

  return s.join("");
}

//eg. format BTC-24SEP21-80000-P
export function getInstrumentName(
  deliveryType,
  optionType,
  quoteCurrency,
  baseCurrency,
  expiry,
  strike
) {
  let OT = optionType === "CALL" ? "C" : "P";
  let DT = deliveryType === "PHYSICAL" ? "P" : "C";
  let months = [
    "JAN",
    "FEB",
    "MAR",
    "APR",
    "MAY",
    "JUN",
    "JUL",
    "AUG",
    "SEP",
    "OCT",
    "NOV",
    "DEC",
  ];
  let date = expiry.Format("d");
  let month = expiry.getMonth();
  let year = expiry.Format("yy");
  return `${DT}-${quoteCurrency}-${baseCurrency}-${date}${months[month]}${year}-${strike}-${OT}`;
}

export function toMoney(num) {
  return parseFloat(num).toLocaleString();
}

export function parseJWTPayload(jwtToken) {
  if (jwtToken && jwtToken.length > 0) {
    let split = jwtToken.split(".");
    if (split.length === 3) {
      let payloadEncode = split[1];
      return JSON.parse(atob(payloadEncode));
    } else {
      return null;
    }
  } else {
    return null;
  }
}

export function isJWTTokenExpired(jwt) {
  let payload = parseJWTPayload(jwt);
  if (payload && payload.exp) {
    return new Date().getTime() > Number(payload.exp * 1000);
  }
  return true;
}

export function isDAppTokenValid(jwt) {
  let payload = parseJWTPayload(jwt);
  if (payload && payload.expiredAt) {
    //expired
    if (new Date().getTime() > Number(payload.expiredAt * 1000)) {
      return false;
    }
    //userId null
    return !(!payload.userId || payload.userId.length === 0);
  }
  return false;
}

export function diffDays(older, newer) {
  let olderDate = new Date(older);
  let newerDate = new Date(newer);
  let differenceInTime = olderDate.getTime() - newerDate.getTime();
  let day = parseInt((differenceInTime / (1000 * 3600 * 24)).toFixed(2));
  return day + 1;
}
