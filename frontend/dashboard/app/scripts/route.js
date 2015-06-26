'use strict';

angular.module("tukdesk").config(function($routeProvider, $locationProvider) {
    $locationProvider.html5Mode({
        enabled: true,
        requireBase: false
    });

    var checkBrand = function($location, $q, globals, api) {
        var d = $q.defer();
        // check brand info
        globals.brandInfoLoad().then(function(data) {
            d.resolve(data);
        }, function(brandErr) {
            api.logResponseObj(brandErr);
            d.reject(brandErr);
            $location.path("/init")
        });
        return d.promise;
    };

    var checkUser = function($location, $q, globals, api) {
        var d = $q.defer();
        // check is agent
        globals.userInfoLoad().then(function(data) {
            if (globals.userIsAgent() === true) {
                d.resolve(data)
            } else {
                d.reject("unauthorized");
            }
        }, function(userErr) {
            api.logResponseObj(userErr);
            d.reject(userErr);
            $location.path("/signin")
        });
        return d.promise;
    };

    var checkAll = function($location, $q, globals, api) {
        var d = $q.defer();
        // check brand first
        checkBrand($location, $q, globals, api).then(function() {
            // check user
            checkUser($location, $q, globals, api).then(function(data) {
                d.resolve(data);
            }, function(userErr) {
                d.reject(userErr)
            })
        }, function(brandErr) {
            d.reject(brandErr);
        });
        return d.promise;
    };

    $routeProvider
        .when("/init", {
            controller: "initCtrl",
            templateUrl: "/views/init.html"
        })
        .when("/signin", {
            controller: "signinCtrl",
            templateUrl: "/views/signin.html",
            resolve: {
                checkBrand: checkBrand
            }
        })
        .when("/tickets", {
            controller: "ticketListCtrl",
            templateUrl: "/views/ticket-list.html",
            resolve: {
                checkAll: checkAll
            }
        })
        .when("/tickets/add", {
            controller: "ticketAddCtrl",
            templateUrl: "/views/ticket-add.html",
            resolve: {
                checkAll: checkAll
            }
        })
        .otherwise({
            redirectTo: "/tickets"
        });
});