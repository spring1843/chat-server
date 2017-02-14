 'use strict';

 // grab the room from the URL
 var room = null;

 // create our webrtc connection
 var webrtc = new SimpleWebRTC({
     // the id/element dom element that will hold "our" video
     localVideoEl: 'localVideo',
     // the id/element dom element that will hold remote videos
     remoteVideosEl: '',
     // immediately ask for camera access
//     autoRequestMedia: true,
     debug: false,
     detectSpeakingEvents: true
 });

 webrtc.on('videoAdded', function (video, peer) {
     console.log('video added', peer);
     var remotes = document.getElementById('remotes');
     if (remotes) {
         var d = document.createElement('div');
         d.className = 'videoContainer';
         d.id = 'container_' + webrtc.getDomId(peer);
         d.appendChild(video);
         var vol = document.createElement('div');
         vol.id = 'volume_' + peer.id;
         vol.className = 'volume_bar';
         video.onclick = function () {
             video.style.width = video.videoWidth + 'px';
             video.style.height = video.videoHeight + 'px';
         };
         d.appendChild(vol);
         remotes.appendChild(d);
     }
 });

 webrtc.on('videoRemoved', function (video, peer) {
     console.log('video removed ', peer);
     var remotes = document.getElementById('remotes');
     var el = document.getElementById('container_' + webrtc.getDomId(peer));
     if (remotes && el) {
         remotes.removeChild(el);
     }
 });

 // Since we use this twice we put it here
 function setRoom(name) {
     $('body').addClass('active');
     room = name;

     webrtc.startLocalVideo();
      // when it's ready, join if we got a room from the URL
      webrtc.on('readyToCall', function () {
          // you can name it anything
          if (room) webrtc.joinRoom(room);
      });

 }

 if (room) {
     setRoom(room);
 }

 var button = $('#screenShareButton'),
     setButton = function (bool) {
         button.text(bool ? 'share screen' : 'stop sharing');
     };
 webrtc.on('localScreenStopped', function () {
     setButton(true);
 });

 setButton(true);

 button.click(function () {
     if (webrtc.getLocalScreen()) {
         webrtc.stopScreenShare();
         setButton(true);
     } else {
         webrtc.shareScreen(function (err) {
             if (err) {
                 setButton(true);
             } else {
                 setButton(false);
             }
         });

     }
 });