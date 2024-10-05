 let cg_params = new Vue({

     el: "#addcg",
     delimiters: ['[[', ']]'],

     data() {
         return {
              
             httpcodes:[],
             httpcodeselect:200,
          }
     },

     computed: {
        
 
      },




     mounted() {
        try {
          
         }
         catch (error) {
            // console.error(error)
 
          }
 
     },
     methods: {

   
 
 

  
 

         donothing(e) {
             e.preventDefault();
             e.stopPropagation();
         },
         loadParams(epid, httpcode, cgparmid) {

             let config = {
                 headers: {
                     "X-CSRF-Token": csrftoken,
                     "Content-Type": "application/json",
                     "Accept": "application/json"
                 }
             }
             let data = {}
             var local = this
             this.showModal = true
             local.processing = true

             axios.post('/query/'+epid+'/'+httpcode+'/'+cgparmid, data, config).then(function (response) {

                console.log(response)

             }).catch(function (error) {

                 // handle error
                 // console.log(error);
             }).then(function () {
                 local.showModal = false
                 local.processing = false
                 handler()
                 // always executed
             });


         },
      
 
    

     }
 });

 