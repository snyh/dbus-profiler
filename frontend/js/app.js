(function() {
    'use strict';

    angular.module('dbus-profiler', ['smart-table', 'ui.router'])
        .config(['$stateProvider',function($stateProvider) {
            var homeState = {
                name: "home",
                url: '/',
                template: '<d-summary></d-summary>'
            }

            var ifcListState = {
                name: 'ifcs',
                url: '/ifcs',
                component: 'dIfcList',
                resolve: {
                    infos: ['dapi', function(dapi) {
                        return dapi.IfcInfos()
                    }]
                }
            }
            var ifcDetailState = {
                name: 'ifcs.detail',
                url: '/{ifcName}',
                component: 'dInterface',
                resolve: {
                    fetchFn: ['$stateParams', 'dapi', function($stateParams, dapi) {
                        return dapi.BuildIfcInfo($stateParams.ifcName)
                    }]
                }
            }
            var methodDetailState = {
                name: 'ifcs.detail.method',
                url: '/{name}?t={type}',
                component: 'dMethod',
                resolve: {
                    fetchFn: ['$stateParams', 'dapi', function($stateParams, dapi) {
                        return dapi.BuildMethodInfo($stateParams.ifcName, $stateParams.type, $stateParams.name)
                    }]
                }
            }

            $stateProvider.state(homeState)
            $stateProvider.state(ifcListState)
            $stateProvider.state(ifcDetailState)
            $stateProvider.state(methodDetailState)
        }]);
})()
