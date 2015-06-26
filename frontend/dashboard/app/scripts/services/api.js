'use strict';

angular.module("tukdesk").factory("api", function($http, $cookies, $resource, $log, apiPrefix, errorCode) {
    var fac = {};
    fac.getAuthHeader = function() {
        return $cookies.get("token") || "";
    };
    fac.setAuthHeader = function(val) {
        $cookies.put("token", val, {"path": "/"});
        fac.setDefaultHeader();
    };
    fac.clearToken = function() {
        $cookies.remove("token", {"path": "/"});
        fac.setDefaultHeader();
    };
    fac.setDefaultHeader = function() {
        $http.defaults.headers.common.Authorization = fac.getAuthHeader();
    };
    fac.apiUrl = function(path) {
        return apiPrefix + path;
    };

    // request
    fac.request = {};
    angular.forEach(['get', 'delete', 'head', 'jsonp'], function(name) {
        fac.request[name] = function(path, config) {
            var url = fac.apiUrl(path);
            config = config || {};
            return $http[name](url, config);
        };
    });
    angular.forEach(['post', 'put'], function(name) {
        fac.request[name] = function(path, data, config) {
            var url = fac.apiUrl(path);
            config = config || {};
            return $http[name](url, data, config);
        };
    });

    // resource
    fac.resBrand = $resource(fac.apiUrl("/brand"), {

    }, {
        get: {method: "GET"},
        update: {method: "PUT"}
    });

    fac.resProfile = $resource(fac.apiUrl("/profile"), {

    }, {
        get: {method: "GET"},
        update: {method: "PUT"}
    });

    fac.resTickets = $resource(fac.apiUrl("/tickets/:ticketId"), {
        ticketId: "@id"
    }, {
        list: {method: "GET"},
        add: {method: "POST"},
        info: {method: "GET"},
        update: {method: "PUT"}
    });

    fac.resComments = $resource(fac.apiUrl("/tickets/:ticketId/comments/:commentId"), {
        commentId: "@id"
    }, {
        list: {method: "GET"},
        add: {method: "POST"},
        update: {method: "PUT"}
    });

    fac.resUsers = $resource(fac.apiUrl("/users/:userId"), {
        userId: "@id"
    }, {
        list: {method: "GET"},
        info: {method: "GET"},
        update: {method: "PUT"}
    });

    fac.resFocus = $resource(fac.apiUrl("/focus/:focusId"), {
        focusId: "@id"
    }, {
        list: {method: "GET"},
        add: {method: "POST"},
        handle: {method: "PUT"}
    });

    // errors
    fac.logHTTPResp = function(data, status, headers) {
        $log.error("Status: " + status + "; "
            + "Code: " + data["error_code"] + "; "
            + "Message: " + data["error_msg"] + "; "
            + "Req-Id: " + headers("X-Req-Id"));
    };

    fac.logResponseObj = function(respObj) {
        fac.logHTTPResp(respObj.data, respObj.status, respObj.headers)
    };

    fac.responseErr = function (cb) {
        return function (resErr) {
            fac.logResponseObj(resErr);
            if (angular.isFunction(cb)) {
                cb(resErr);
            }
        }
    };

    fac.httpErr = function (cb) {
        return function (data, status, headers, config) {
            fac.logHTTPResp(data, status, headers);
            if (angular.isFunction(cb)) {
                cb(data, status, headers, config);
            }
        }
    };

    fac.isErrorType = function(code, typ) {
        return errorCode[typ] === code;
    };

    return fac;
});
