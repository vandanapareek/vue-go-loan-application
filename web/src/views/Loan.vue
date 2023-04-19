<template>
    <div id="overlay" :class="loading == true ? 'overlay_active' : 'overlay_inactive'"></div>
    <v-card class="text-center">
      <v-card-text>
        <v-row>
          <v-col cols="2" class="header-img">
                    <img class="avatar_img" :src="getPic()" alt="Avatar" style="width:50px;height: 50px;"/>
                    <span class="avatar_title">Hi, {{ username }}</span>
          </v-col>
          <v-col cols="8" class="header-title">
            <h1 class="display-1">            
              <img src="../assets/wallet_icon.png" width="40" height="40" class="header-img"/>
              Loan Calculator 
            </h1>

          </v-col>
          <v-col cols="2" class="header-btn">
            <v-btn class="btn-warning" color="warning" type="warning" @click="logout()">
              Log Out
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
 

      <div class="wrapper" >
        
        <div class="card-form">
          <div class="card-form__inner">
            <div class="card-input">
              <label for="businessName" class="card-input__label">Business Name</label>
              <input type="text" id="businessName" class="card-input__input" v-model="businessName" autocomplete="off">
            </div>

            <div class="card-input">
              <label for="estYear" class="card-input__label">Year Established</label>
              <input type="text" id="estYear" class="card-input__input" v-model="estYear" autocomplete="off" @keypress="isNumber($event)">
            </div>

            <div class="card-input">
              <label for="loanAmount" class="card-input__label">Loan Amount (SGD)</label>
              <input type="text" id="loanAmount" class="card-input__input" v-model="loanAmount" autocomplete="off" @keypress="isNumber($event)">
            </div>

            <div class="card-input">
              <label for="provider" class="card-input__label">Provider</label>
                  <select class="card-input__input -select" id="provider" v-model="provider"  data-ref="cardDate">
                    <option value="" disabled selected>Choose Provider</option>
                    <option v-for="item in providers" :key='item.slug' v-bind:value="item.slug">
                      {{item.name}}
                    </option>
                  </select>
            </div>
            <div v-if="loading" id="loading">
              <img id="loading-image" src="../assets/send.gif" alt="Loading..." />
            </div>
            <button class="card-form__button" @click="getBalanceSheet">
              Request Balance Sheet
            </button>
          </div>
        </div> 
      </div>

      <div class="d-flex align-center flex-column" v-if="balanceSheet.sheet">
        <v-card
        width="800"
        >
        <v-card-text >
            <v-row no-gutters >
              <v-col cols="6" class="pl-8 pt-6"><h2>Balance Sheet</h2></v-col>
              <v-col cols="6" class="text-right pr-8 pt-4">
                <v-btn class="btn-success" color="success" type="success" @click="getLoanDetails()">
                  Loan Details
              </v-btn>
              </v-col>
            </v-row>
          </v-card-text>

          <v-card-text v-if="preAssessment != -1">
            <v-row no-gutters class="bg-yellow" >
              <v-col
                cols="6" 
              >
                <v-sheet class="pa-2 ma-2">
                  Loan Permitted (SGD): {{ loanPermit }}
                </v-sheet>
              </v-col>

              <v-col
                cols="6"
              >
                <v-sheet class="pa-2 ma-2">
                  Pre Assessment: {{ preAssessment }}%
                </v-sheet>
              </v-col>
            </v-row>
          </v-card-text>

          <BalanceSheet :balanceSheet="balanceSheet"/> 
        </v-card>
      </div>

    <v-snackbar v-model="snackbar"  :timeout="timeout" :color="color" location="top">
          {{ snackbarText }}
    </v-snackbar>
</template>

<script>
import axios from 'axios';
import BalanceSheet  from '../components/BalanceSheet.vue';
import {VDataTable} from 'vuetify/labs/VDataTable'

export default {
   
  el: "#app",
  components: {
    BalanceSheet,
    VDataTable,
  },

  data() {
    
    return {
      username:localStorage.getItem('user_name'),
      providers:[],
      businessName:"",
      estYear:"",
      loanAmount:"",
      balanceSheet:{},
      loading: false,
      snackbar: false,
      snackbarText: "",
      timeout: 3000,
      color: "",
      search: '',
      loanPermit: 0,
      preAssessment: -1
    };
  },
  mounted() {
    const config = {
      headers: { Authorization: `Bearer `+localStorage.getItem("loan_token") }
    };
    axios.get("/providers",config).then(response => {
           this.providers = response.data;
           console.log(this.providers)
        }).catch((error) => {
          console.log(error)
          if((error.code == "ERR_NETWORK") || (error.response && error.response.status == 401)) {
            this.logout()
          }
        })
  },
  methods: {
     getBalanceSheet() {
      this.balanceSheet = {}
      this.preAssessment = -1
      this.loanPermit = 0
      var current_year=new Date().getFullYear();
      if (this.estYear.length != 4 || (this.estYear < 1920) || (this.estYear > current_year)) {
        this.snackbar = true;
        this.color = 'red';
        this.snackbarText = "Please enter valid year"
      } else {
        var formData = {
            businessName : this.businessName,
            estYear : parseInt(this.estYear),
            provider: this.provider,
            loanAmount: parseInt(this.loanAmount)
        }
        const config = {
        headers: { Authorization: `Bearer `+localStorage.getItem("loan_token"), 'Content-Type': 'application/json'  }
        };
         axios.post("/balance-sheet", JSON.stringify(formData) ,config).then(response => {
              this.isActive = !this.isActive;
              console.log(response.data)
              this.balanceSheet = response.data
              this.$nextTick(() => {
                this.bottom = document.body.scrollHeight
                window.scrollTo(0, this.bottom);
              })
                    
        }).catch((error) => {
            console.log(error)  
            if(error.response && error.response.data.Code == 422) {
              this.snackbar = true;
              this.color = 'red';
              this.snackbarText = error.response.data.Message
            } 
            if((error.code == "ERR_NETWORK") || (error.response &&  error.response.data.Code == 401)) {
              this.logout()
            } 
          })  
      }
    },
    getLoanDetails() {
      var formData = {
            businessName : this.businessName,
            estYear : parseInt(this.estYear),
            loanAmount: parseInt(this.balanceSheet.loanAmount),
            totalProfit : parseInt(this.balanceSheet.totalProfit),
            avgAssets: parseInt(this.balanceSheet.avgAssets),
        }
        const config = {
        headers: { Authorization: `Bearer `+localStorage.getItem("loan_token"), 'Content-Type': 'application/json'  }
        };
        axios.post("/calculate-loan", JSON.stringify(formData) ,config).then(response => {
              console.log(response.data)
              this.loanPermit = response.data.loanPermitVal
              this.preAssessment = response.data.preAssessment
        }).catch((error) => {
            console.log(error)  
            if(error.response && error.response.data.Code == 422) {
              this.snackbar = true;
              this.color = 'red';
              this.snackbarText = error.response.data.Message
            } 
            if((error.code == "ERR_NETWORK") || (error.response &&  error.response.data.Code == 401)) {
              this.logout()
            } 
          })  
    },
    getPic() {
       return 'src/assets/' + this.username + '.jpeg';
    },
    logout() {
      localStorage.clear();
      this.$router.go(0);
    },
    isNumber: function(evt) {
      evt = (evt) ? evt : window.event;
      var charCode = (evt.which) ? evt.which : evt.keyCode;
      if ((charCode > 31 && (charCode < 48 || charCode > 57) || (charCode > 31 && (charCode < 48 || charCode > 57 )) ) ) {
        evt.preventDefault();;
      } else {
        return true;
      }
    },
    yearValidation(year) {
      var text = /^[0-9]+$/;
        if (year != 0) {
            if ((year != "") && (!text.test(year))) {
                this.snackbar = true;
                this.color = 'red';
                this.snackbarText ="Please Enter Numeric Values Only"
                return false;
            }

            if (year.length != 4) {
                this.snackbar = true;
                this.color = 'red';
                this.snackbarText ="Year is not proper. Please check"
                return false;
            }
            var current_year=new Date().getFullYear();
            if((year < 1920) || (year > current_year))
                {
                  this.snackbar = true;
                  this.color = 'red';
                  this.snackbarText ="Year should be in range 1920 to current year"
                  return false;
                }
            return true;
        } 
    }
  }
};
</script>
