var API_KEY = 'pk.eyJ1IjoibWphbmdlciIsImEiOiJjazN6NHZlNHkwMjZiM2tudzRpN3FyNzc0In0.avUDs9ardvviib8L8HsMSA';
var tileLayer = L.tileLayer('https://api.mapbox.com/styles/v1/{id}/tiles/{z}/{x}/{y}?access_token={accessToken}', {
    attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, <a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
    maxZoom: 18,
    //id: 'mapbox/streets-v11',
    id: 'mapbox/light-v10',
    tileSize: 512,
    zoomOffset: -1,
    accessToken: API_KEY
});

var RecentMemories = RecentMemories || (function() {
    var infoWindow, map, bounds, markers, prevMarker, group, ui;    
        
    function init() {
        map = L.map('map').setView([38.677811, -90.419197], 13);

        tileLayer.addTo(map);
        
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
        var marker = L.marker([latitude, longitude]);
        marker.bindPopup('<a href="/memories/' + memoryId + '">' + memoryTitle + '</a>');
        marker.addTo(map);

        markers[memoryId] = marker;

        if (!bounds) {
            bounds = new L.LatLngBounds(); 
        }

        bounds.extend(marker.getLatLng());

        map.fitBounds(bounds);
    }

    function showInfoWindowForMemory(memoryId) {
        var marker = markers[memoryId];
        marker.openPopup();
    }

    return {
        showInfoWindowForMemory: showInfoWindowForMemory,
        addMemory: addMemory, 
        init: init
    }
})();

var AddMemory = AddMemory || (function() {
    var map, marker;

    function mapClicked(e) {
        if (!marker) {
            marker = L.marker(e.latlng);
            marker.addTo(map)
        }
        
        marker.setLatLng(e.latlng);
        
        $("#latitude").val(e.latlng.lat);
        $("#longitude").val(e.latlng.lng);
    }

    function findAddress() {
        var latitude = $("#latitude").val();
        var longitude = $("#longitude").val();

        if (latitude && longitude) {
            marker = L.marker([latitude, longitude]);
            marker.addTo(map)
        }
    }

    function init() {
        map = L.map('address_search_map').setView([38.6270, -90.1994], 4);

        tileLayer.addTo(map);

        map.on('click', mapClicked);
    }

    function verifyForm() {
        var latitude = $("#latitude").val();
        var longitude = $("#longitude").val();

        if (!latitude || !longitude) {
            alert("Please pick a location on the map.");
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
        var map = L.map('memory_map').setView([lat, lng], 13);

        tileLayer.addTo(map);
        
        L.marker([lat, lng]).addTo(map);
    }

    return {
        showMap: showMap
    }
})();