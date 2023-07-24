var Filters = {
  byteFormat: function(size) {
    var result;
    switch (true) {
      case size === null || size === "" || isNaN(size):
        result = "-";
        break;

      case size >= 0 && size < 1024:
        result = size + " B";
        break;

      case size >= 1024 && size < Math.pow(1024, 2):
        result = Math.round(size / 1024 * 100) / 100 + " KB";
        break;

      case size >= Math.pow(1024, 2) && size < Math.pow(1024, 3):
        result = Math.round(size / Math.pow(1024, 2) * 100) / 100 + " MB";
        break;

      case size >= Math.pow(1024, 3) && size < Math.pow(1024, 4):
        result = Math.round(size / Math.pow(1024, 3) * 100) / 100 + " GB";
        break;

      default:
        result = Math.round(size / Math.pow(1024, 4) * 100) / 100 + " TB";
    }
    return result;
  },
  secondsInReadable: function(value) {
    if (!value) return "-";

    var d, h, m, s;
    s = value;
    m = Math.floor(s / 60);
    s = s % 60;
    h = Math.floor(m / 60);
    m = m % 60;
    d = Math.floor(h / 24);
    h = h % 24;

    if (d == 0) {
      var time = h + "h " + m + "m " + s + "s";
    } else {
      var time = d + "d " + h + "h " + m + "m";
    }

    return time;
  },
  formatDate: function(value, hours = true) {
    var moment = require("moment/moment.js");
    if (+new Date(value) > 0)
      return moment(value).format("DD MMMM YYYY" + (hours ? ", HH:mm" : ""));
    else return "-";
  },
  dateFromNow: function(value) {
    var moment = require("moment/moment.js");
    var date = new Date(value);
    if (+date > 0) return moment(date).fromNow();
    else return "-";
  }
};

export default Filters;
