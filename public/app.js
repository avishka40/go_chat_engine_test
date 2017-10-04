
new Vue({
  el:'#app',

  data:{

    ws: null,//websocket
    newMsg:'',//holds a message
    chatContent:'',//running list of chat messags diplaed on screen
    email:null,//username
    username:null,
    joined:false

  },
  //something similar to a constrcutor ,whenever a vue instance is created this proceeds with the function specified in the created:attribute
  created:function(){
    var self =this;
    this.ws = new WebSocket('ws://'+window.location.host + '/ws')
    this.ws.addEventListener('message',function(e){
      var msg =JSON.parse(e.data);
      self.chatContent += '<div class="chip">'
        + '<img src="' +self.gravatarURL(msg.email)+'">'//avatar
        +msg.username
      +'</div>'
      +emojione.toImage(msg.message) + '</br>';//parse emojis

      var element =document.getElementById('chat-messages');
      element.scrollTop = element.scrollHeight;

    });
  },
  methods:{
    send: function(){
      if (this.newMsg!=''){
        this.ws.send(
          JSON.stringify({
            email:this.email,
            username:this.username,
            message:$('<p>').html(this.newMsg).text()//strip out html
          })
        );
        this.newMsg ='';
      }
    },
    join:function(){
      if(!this.email){
        Materialize.toast('You must enter an email',2000);
        return
      }
      if(!this.username){
        Materialize.toast('You must enter an username',2000);
        return
      }
      this.email=$('<p>').html(this.email).text();
      this.username =$('<p>').html(this.username).text();
      this.joined = true;
    },
    gravatarURL:function(email){
      return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
    }
  }



});
