new Vue({
  el: '#app',
  
  data: {
    ws: null, // Our websocket
    newMsg: '', // Holds new messages to be sent to the server
    chatContent: '', // A running list of chat messages displayed on the screen
    email: null, // Email address used for grabbing an avatar
    nickname: null, // Our nickname
    joined: false // True if email and nickname have been filled in
  },
  
  created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
      var msg = JSON.parse(e.data);
      self.chatContent += '<div class="chip">'
      + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
      + msg.nickname
      + '</div>'
      + emojione.toImage(msg.text) + '<br/>'; // Parse emojis
      
      var element = document.getElementById('chat-messages');
      element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
    });
  },
  
  methods: {
    send: function () {
      if (this.newMsg != '') {
        this.ws.send(
          JSON.stringify({
            email: this.email,
            nickname: this.nickname,
            text: $('<p>').html(this.newMsg).text(), // Strip out html
            timestamp: new Date()
          }
        ));
        this.newMsg = ''; // Reset newMsg
      }
    },
    
    join: function () {
      if (!this.email) {
        Materialize.toast('You must enter an email', 2000);
        return
      }
      if (!this.nickname) {
        Materialize.toast('You must choose a nickname', 2000);
        return
      }
      this.email = $('<p>').html(this.email).text();
      this.nickname = $('<p>').html(this.nickname).text();
      this.joined = true;
    },
    
    gravatarURL: function(email) {
      return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
    }
  }
});