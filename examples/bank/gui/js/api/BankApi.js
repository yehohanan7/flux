import axios from 'axios';

class BankApi {
  static getAllEvents(callback) {
    console.log('get Events');
    /*
    return axios.get('/bank/events').then(response => {
      callback(response.data)
    }).catch(error => {
      throw(error);
    });
    */
  }

  static currentBalance(callback) {
    console.log('current balance');
  }

  static credit(amount, callback) {
    console.log('credit req');
    console.log(amount);
  }

  static debit(amount, callback) {
    console.log('debit req');
    console.log(amount);
  }
}

export default BankApi;