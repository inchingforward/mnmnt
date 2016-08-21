var MAP_STYLES = [{"featureType":"landscape","stylers":[{"hue":"#FFBB00"},{"saturation":43.400000000000006},{"lightness":37.599999999999994},{"gamma":1}]},{"featureType":"road.highway","stylers":[{"hue":"#FFC200"},{"saturation":-61.8},{"lightness":45.599999999999994},{"gamma":1}]},{"featureType":"road.arterial","stylers":[{"hue":"#FF0300"},{"saturation":-100},{"lightness":51.19999999999999},{"gamma":1}]},{"featureType":"road.local","stylers":[{"hue":"#FF0300"},{"saturation":-100},{"lightness":52},{"gamma":1}]},{"featureType":"water","stylers":[{"hue":"#0078FF"},{"saturation":-13.200000000000003},{"lightness":2.4000000000000057},{"gamma":1}]},{"featureType":"poi","stylers":[{"hue":"#00FF6A"},{"saturation":-1.0989010989011234},{"lightness":11.200000000000017},{"gamma":1}]}];

var RecentMemories = RecentMemories || (function() {
    var infoWindow, map, bounds, markerAndContent;
    
    function init() {
        var mapOptions = {
            center: { lat: 38.677811 , lng:  -90.419197 },
            zoom: 13,
            zoomControl: true,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false,
            styles: MAP_STYLES
        };

        infoWindow = new google.maps.InfoWindow();
        map = new google.maps.Map(document.getElementById("map"), mapOptions);
        bounds = new google.maps.LatLngBounds();
        markerAndContent = {};
    }

    function addMemory(latitude, longitude, memoryId, memoryTitle) {
        var latLng = new google.maps.LatLng(latitude, longitude);

        bounds.extend(latLng);

        var marker = new google.maps.Marker({
            position: latLng,
            map: map,
            title: memoryTitle
        });

        marker.addListener("mouseover", function() {
            var that = this;
            showInfoWindowForId(memoryId);
        });

        cacheMarkerAndContent(memoryId, marker, memoryTitle);

        map.fitBounds(bounds);        
    }

    function cacheMarkerAndContent(memoryId, marker, content) {
        markerAndContent[memoryId] = {
            marker: marker,
            content: content
        }
    }

    function showInfoWindowForId(memoryId) {
        var mc = markerAndContent[memoryId];
        var content = '<a href="/memories/' + memoryId + '">' + mc.content + '</a>'

        infoWindow.setContent(content);
        infoWindow.open(map, mc.marker);
    }

    return {
        showInfoWindowForId: showInfoWindowForId,
        addMemory: addMemory, 
        init: init
    }
})();

var AddMemory = AddMemory || (function() {
    var map, marker;

    var delay = function() {
        var timeout = 0;
        return function(callback, ms) {
            clearTimeout(timeout);
            timeout = setTimeout(callback, ms);
        };
    }();

    function init() {
        var mapOptions = {
            center: { lat: 38.6272222, lng: -90.1977778},
            zoom: 13,
            zoomControl: true,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false,
            styles: MAP_STYLES
        };

        map = new google.maps.Map(document.getElementById('address_search_map'), mapOptions);

        $("#address_text").keyup(function() {
            delay(findAddress, 1500);
        });

        $("#address_text").keypress(function (e) {
            if ((e.which && e.which == 13) || (e.keyCode && e.keyCode == 13)) {
                findAddress();
            }
        });
    }

    function findAddress() {
        var addressText = $("#address_text").val();
        if (!addressText) {
            alert("Please enter an address to search for.");
            return;
        }

        if (marker) {
            marker.setMap(null);
        }

        var latitude = $("#latitude");
        var longitude = $("#longitude");
        var submitButton = $("#submit_button");
        var geocoder = new google.maps.Geocoder();

        $("#map_spinner").show();

        geocoder.geocode({address: addressText}, function(results, status) {
            $("#map_spinner").hide();

            if (status == google.maps.GeocoderStatus.OK) {
                map.setCenter(results[0].geometry.location);

                marker = new google.maps.Marker({
                    map: map,
                    position: results[0].geometry.location
                });

                latitude.val(results[0].geometry.location.lat());
                longitude.val(results[0].geometry.location.lng());

                submitButton.attr("disabled", false);
            } else {
                latitude.val("");
                longitude.val("");

                submitButton.attr("disabled", "disabled");
            }
        });
    }

    function verifyForm() {
        var latitude = $("#latitude").val();
        var longitude = $("#longitude").val();

        if (!latitude || !longitude) {
            alert("Please search and find an address.");
            $("#address_text").focus();
            return false;
        }

        var titleField = $("#title");
        if (!titleField.val()) {
            alert("Please enter a title.");
            titleField.focus();
            return false;
        }

        var detailsField = $("#details");
        if (!detailsField.val()) {
            alert("Please enter some details.");
            detailsField.focus();
            return false;
        }

        return true;
    }

    return {
        init: init,
        findAddress: findAddress
    }
})();

var EditMemory = AddMemory || (function() {
    var map, marker;

    function init(latitude, longitude) {
        var position = { lat: latitude, lng: longitude };

        var mapOptions = {
            center: position,
            zoom: 13,
            zoomControl: true,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false,
            styles: MAP_STYLES
        };

        map = new google.maps.Map(document.getElementById('edit_memory_map'), mapOptions);

        marker = new google.maps.Marker({
            map: map,
            position: position
        });
    }

    function verifyForm() {
        var detailsField = $("#details");
        if (!detailsField.val()) {
            alert("Please enter some details.");
            detailsField.focus();
            return false;
        }

        return true;
    }

    return {
        init: init,
    }
})();

var MemoryDetails = MemoryDetails || (function() {
    function showMap(lat, lng, title) {
        var mapOptions = {
            center: { lat: lat, lng: lng},
            zoom: 13,
            zoomControl: true,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false,
            styles: MAP_STYLES
        };

        var map = new google.maps.Map(document.getElementById("memory_map"), mapOptions);

        new google.maps.Marker({
            position: {lat: lat, lng: lng}, 
            map: map, 
            title: title
        });
    }

    return {
        showMap: showMap
    }
})();