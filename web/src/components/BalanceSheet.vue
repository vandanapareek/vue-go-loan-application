<template>
  <v-sheet class="pa-12" rounded>
    <v-card>
      <v-data-table
      :fixedFooter="false"
      :headers="headers"
      :items="balanceSheet.sheet"
      class="elevation-1"
     ></v-data-table>

     <v-card-text>
        <v-row no-gutters class="bg-green" >
          <v-col
            cols="6"
          >
            <v-sheet class="pa-2 ma-2">
              Total Profit (SGD): {{ balanceSheet.totalProfit }}
            </v-sheet>
          </v-col>

          <v-col
            cols="6"
          >
            <v-sheet class="pa-2 ma-2">
              Avg Assets (SGD): {{ balanceSheet.avgAssets }}
            </v-sheet>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-sheet>
</template>

<script>
import { VDataTable } from 'vuetify/labs/VDataTable'

export default {
  components: {
    VDataTable
  },
  props: {
    balanceSheet: Object
  },
  data() {
    return {
      headers: [
      { title: 'Year', key: 'year' },
      { title: 'Month', key: 'month' },
      { title: 'Profit/Loss (SGD)', key: 'profitOrLoss' },
      { title: 'Assets (SGD)', key: 'assetsValue' },
    ],

    }
  },
  name: "BlanceSheet"
};
</script>