'use strict';

angular.module("tukdesk").run(function($rootScope, api, globals) {
    console.log("running");

    api.setDefaultHeader();
    $rootScope.globals = globals;
});
