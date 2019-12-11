var API_KEY = 'VH-NAclLOgOGa9AKWvrMO9ttZHkTDKL71nxUf2jL7bM';

var RecentMemories = RecentMemories || (function() {
    var infoWindow, map, bounds, markers, prevMarker, group, ui;    
        
    function init() {
        var center = { lat: 38.677811, lng: -90.419197 };

        var platform = new H.service.Platform({
            'apikey': API_KEY
        });
    
        var defaultLayers = platform.createDefaultLayers();
    
        map = new H.Map(
            document.getElementById('map'),
            defaultLayers.vector.normal.map, 
            {
                zoom: 13,
                center: center,
                pixelRatio: window.devicePixelRatio || 1,
                padding: {top: 15, right: 15, bottom: 15, left: 15}
            }
        );

        var behavior = new H.mapevents.Behavior(new H.mapevents.MapEvents(map));
        ui = H.ui.UI.createDefault(map, defaultLayers);

        group = new H.map.Group();

        map.addObject(group);

        group.addEventListener('tap', function (evt) {
            showMarkerBubble(evt.target);
        }, false);

        markers = {};
    }

    function showMarkerBubble(marker) {
        var bubbles = ui.getBubbles();
        for (var i=0; i < bubbles.length; i++) {
            ui.removeBubble(bubbles[i]);
        }

        var bubble =  new H.ui.InfoBubble(marker.getGeometry(), {
            content: marker.getData()
        });

        ui.addBubble(bubble);
    }

    function addMemory(latitude, longitude, memoryId, memoryTitle) {
        var marker = new H.map.Marker({lat: latitude, lng: longitude});
        marker.setData('<a href="/memories/' + memoryId + '">' + memoryTitle + '</a>');
        group.addObject(marker);
        
        map.getViewModel().setLookAtData({
             bounds: group.getBoundingBox()
        });

        markers[memoryId] = marker;
    }

    function showInfoWindowForMemory(memoryId) {
        var marker = markers[memoryId];
        showMarkerBubble(marker);
    }

    return {
        showInfoWindowForMemory: showInfoWindowForMemory,
        addMemory: addMemory, 
        init: init
    }
})();

var AddMemory = AddMemory || (function() {
    var map, marker, geocoder;

    var delay = function() {
        var timeout = 0;
        return function(callback, ms) {
            clearTimeout(timeout);
            timeout = setTimeout(callback, ms);
        };
    }();

    function init() {
        map = new mapboxgl.Map({
            container: 'address_search_map',
            style: 'mapbox://styles/mapbox/light-v10',
            center: [-90.1994, 38.6270],
            zoom: 13
        });

        geocoder = new MapboxGeocoder({
            accessToken: mapboxgl.accessToken,
            mapboxgl: mapboxgl,
            placeholder: "1. Search for a place or address" 
        });

        //document.getElementById('geocoder').appendChild(geocoder.onAdd(map));
        
        $("#address_text").keyup(function() {
            delay(findAddress, 1500);
        });

        $("#address_text").keypress(function (e) {
            if ((e.which && e.which == 13) || (e.keyCode && e.keyCode == 13)) {
                findAddress();
            }
        });

        // geocoder = new MapboxGeocoder({ accessToken: mapboxgl.accessToken });
        geocoder.on('result', function(result) {
            console.log(result);
            console.log(result.result.center);
        });

        console.log(geocoder);
        console.log(geocoder.query);
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

        console.log("gonna query...")
        //geocoder.query(addressText);
        //geocoder.inputString = addressText;
        geocoder._geocode(addressText);
        console.log("queried");

        // var latitude = $("#latitude");
        // var longitude = $("#longitude");
        // var submitButton = $("#submit_button");
        // var geocoder = new google.maps.Geocoder();

        // $("#map_spinner").show();


        // geocoder.geocode({address: addressText}, function(results, status) {
        //     $("#map_spinner").hide();

        //     if (status == google.maps.GeocoderStatus.OK) {
        //         map.setCenter(results[0].geometry.location);

        //         marker = new google.maps.Marker({
        //             map: map,
        //             position: results[0].geometry.location
        //         });

        //         latitude.val(results[0].geometry.location.lat());
        //         longitude.val(results[0].geometry.location.lng());

        //         submitButton.attr("disabled", false);
        //     } else {
        //         latitude.val("");
        //         longitude.val("");

        //         submitButton.attr("disabled", "disabled");
        //     }
        // });
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
        var lngLat = [lng, lat];

        var map = new mapboxgl.Map({
            container: 'edit_memory_map',
            style: 'mapbox://styles/mapbox/light-v10',
            zoom: 13,
            center: lngLat
        });
        
        var marker = new mapboxgl.Marker()
            .setLngLat(lngLat)
            .addTo(map);
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
        var latLng = {lat: lat, lng: lng};
        var platform = new H.service.Platform({
            'apikey': API_KEY
        });
    
        var defaultLayers = platform.createDefaultLayers();
    
        map = new H.Map(
            document.getElementById('memory_map'),
            defaultLayers.vector.normal.map, 
            {
                zoom: 13,
                center: latLng,
                pixelRatio: window.devicePixelRatio || 1,
                padding: {top: 15, right: 15, bottom: 15, left: 15}
            }
        );

        var behavior = new H.mapevents.Behavior(new H.mapevents.MapEvents(map));
        ui = H.ui.UI.createDefault(map, defaultLayers);
        
        var marker = new H.map.Marker(latLng);
        map.addObject(marker);
    }

    return {
        showMap: showMap
    }
})();