import { toMoney } from "@/util/utils";
import { SideString } from "@/util/option";

export function instrumentCardTitle(instrument, side, locale) {
  let i = instrument;
  if (!i) {
    return;
  }
  let sideString = SideString(side, instrument.optionType, locale);
  return `$${toMoney(i.strikePrice)} ${sideString} ${i.baseCurrency}`;
}

export function parseInstrumentName(name) {
  let idxDeliveryType = 1;
  let idxQuoteCurrency = 2;
  let idxBaseCurrency = 3;
  let idxDay = 4;
  let idxMonth = 5;
  let idxYear = 6;
  let idxStrikePrice = 7;
  let idxOptionType = 8;

  let instrument = {};
  let monthMap = {
    JAN: "01",
    FEB: "02",
    MAR: "03",
    APR: "04",
    MAY: "05",
    JUN: "06",
    JUL: "07",
    AUG: "08",
    SEP: "09",
    OCT: "10",
    NOV: "11",
    DEC: "12",
  };
  let reg =
    /^([CP])-(pUSD|USDT)-(BTC|XIN|ETH)-([\d]{1,2})(JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)([\d]{2})-([\d]+)-([CP])$/;
  let match = name.match(reg);
  let day = match[idxDay];
  if (day.length === 1) {
    day = "0" + day;
  }
  let timeString = `20${match[idxYear]}-${match[idxMonth]}-${day}T16:00:00+08:00`;
  let expirationDate = new Date(timeString);
  let strikePrice = match[idxStrikePrice];
  let deliveryType = "";
  if (match[idxDeliveryType] === "P") {
    deliveryType = "PHYSICAL";
  }
  if (match[idxDeliveryType] === "C") {
    deliveryType = "CASH";
  }
  let optionType = "";
  if (match[idxOptionType] === "P") {
    optionType = "PUT";
  }
  if (match[idxOptionType] === "C") {
    optionType = "CALL";
  }
  if (match.length === 9) {
    instrument = {
      deliveryType: deliveryType,
      quoteCurrency: match[idxQuoteCurrency],
      baseCurrency: match[idxBaseCurrency],
      expirationDate: expirationDate,
      expirationTimestamp: expirationDate.getTime(),
      strikePrice: strikePrice,
      optionType: optionType,
    };
  }
  return instrument;
}

// The deribit rule, tomorrow is the day after tomorrow, the day after tomorrow +7 the day after tomorrow +14,
// (the day after tomorrow +14) the last Friday of the next month, (the day after tomorrow +14)
export function getBidExpiryDates() {
  let dates = [];
  let firstDay = new Date();
  firstDay.setDate(firstDay.getDate() + 1);
  dates.push(firstDay.Format("yyyy/MM/dd"));
  let secondDay = new Date(firstDay);
  secondDay.setDate(secondDay.getDate() + 1);
  let lastWeeksDay = new Date();
  let weeksDay = Array.from({ length: 3 }, (item, index) => {
    if (index === 0) {
      return secondDay.Format("yyyy/MM/dd");
    } else {
      let date = new Date(secondDay);
      date.setDate(secondDay.getDate() + 7 * index);
      lastWeeksDay = date;
      return date.Format("yyyy/MM/dd");
    }
  });
  let result = dates.concat(weeksDay);

  let diffList = [1, 2, 4, 7, 10];
  let year = lastWeeksDay.getFullYear();
  let isYearIncr = false;
  for (let i in diffList) {
    let month = lastWeeksDay.getMonth() + diffList[i];
    if (month > 11 && !isYearIncr) {
      year++;
      isYearIncr = true;
    }
    let lastFriday = lastFridayOfMonth(year, month);
    result = result.concat(lastFriday.Format("yyyy/MM/dd"));
  }

  return result;
}

export function toExpiryDate(yyMMdd) {
  let strings = yyMMdd.split("/");
  if (strings.length > 1) {
    let year = Number(strings[0]);
    let monthIdx = Number(strings[1]) - 1;
    let day = Number(strings[2]);
    let date = new Date(year, monthIdx, day);
    date.setHours(16, 0, 0, 0);
    return date;
  } else {
    return null;
  }
}

export function leftSettlementPeriod(expirationDate) {
  if (expirationDate) {
    let now = new Date();
    let nowTs = now.getTime();
    let expiry = new Date(expirationDate);
    let expiryTs = new Date(expirationDate).getTime();
    let differenceInTime = expiryTs - nowTs;
    let day = parseInt(differenceInTime / (1000 * 3600 * 24));
    if (day === 0) {
      let hour = parseInt(differenceInTime / (1000 * 3600));
      return `${hour} hours`;
    }
    return `${day + 1} days`;
  } else {
    return "0 days";
  }
}

function lastFridayOfMonth(year, month) {
  "use strict";
  var lastDay = new Date(year, month + 1, 0);
  if (lastDay.getDay() < 6) {
    lastDay.setDate(lastDay.getDate() - 7);
  }
  lastDay.setDate(lastDay.getDate() - (lastDay.getDay() - 6));
  return lastDay;
}
