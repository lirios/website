var app = angular.module("LiriCLA", ["ngRoute"], function($routeProvider) {
    $routeProvider
    .when("/info", {
        title: "info",
        templateUrl: "/static/cla/templates/info.html"
    })
    .otherwise({
        redirectTo: '/info'
    });
});

app.run(['$location', '$rootScope', function($location, $rootScope) {
    $rootScope.$on("$routeChangeSuccess", function(event, current, previous) {
        if (current.$$route) {
            $rootScope.title = current.$$route.title;
        }
    });
}]);

app.controller("MainCtl", function($scope, $timeout, Run) {
});
