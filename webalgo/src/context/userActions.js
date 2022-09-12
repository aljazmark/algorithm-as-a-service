import axios from 'axios';
import {isEmpty} from "lodash"
import Cookies from 'js-cookie'
//const API_URL = 'https://algo-algo.herokuapp.com';
const API_URL = 'http://localhost:8090'
export async function login(dispatch,loginPayload) {
        dispatch({ type: 'LOGIN_REQUEST'});
        var data = {};
        await axios.post(`${API_URL}/user/login`,JSON.stringify(loginPayload)
        ).then(res => {
            data = res
            if (data.data.userID) {
                Cookies.set('user', JSON.stringify(data.data), { expires: 1 })
                dispatch({ type: 'LOGIN_OK', payload: data.data });
            }
            if(data.data.message){
                dispatch({ type: 'LOGIN_ERROR', error: data.data.message });
                
            }else{
                dispatch({ type: 'LOGIN_ERROR', error: "Error logging in" });
            }
        })
        .catch(err =>{
            data = err.response
            dispatch({ type: 'LOGIN_ERROR', error: "Error logging in" });
        });            
    return data || "empty login return";
};

export async function register(dispatch,registerPayload) {
        var data = {};
        await axios.post(`${API_URL}/user`,JSON.stringify(registerPayload)
        ).then(res => {
            data = res
            if (data.data.userID) {
                Cookies.set('user', JSON.stringify(data.data), { expires: 1 })
                dispatch({ type: 'LOGIN_OK', payload: data.data });
            }
            if(data.data.message){
                dispatch({ type: 'LOGIN_ERROR', error: data.data.message });
                
            }else{
                dispatch({ type: 'LOGIN_ERROR', error: "Error logging in" });
            }
        })
        .catch(err =>{
            data = err.response
            dispatch({ type: 'LOGIN_ERROR', error: "Error logging in" });
        });
    
        
    return data || "empty registration return";
}

export async function getUser(dispatch,payload){
    var data = {};
    let config = {
        headers: {'Authorization': `Bearer ${payload.token}`}
    }
        await axios.get(`${API_URL}/user/`+payload.userID,config).then(res => {
            data = res
            if (data.data) {
                dispatch({ type: 'ADD_EMAIL', payload: data.data });
            }
        })
        .catch(err =>{
            if(err.response.status===401){
                alert("Your session has expired, please log in.")
                logout(dispatch)
            }
        });
    return data.data
}
export async function updateUser(dispatch,payload){
    var data = {message:"User update failed"};
    let config = {
        headers: {'Authorization': `Bearer ${payload.token}`}
    }
    var input = JSON.stringify(payload.updateBody) 
        await axios.put(`${API_URL}/user/`+payload.userID,input,config).then(res => {
            if (res.status===200) {
                logout(dispatch)
            }
        })
        .catch(err =>{
            if(err.response.status===401){
                alert("Your session has expired, please log in.")
                logout(dispatch)
            }else if(err.response.status===400){
                data = err.response.data
            }else{
                data = {message:"User update failed"}
            }
        });
    return data
}

export async function getAlgorithms(dispatch,payload){
    var data = {};
        await axios.get(`${API_URL}/help/algorithms`).then(res => {
            data = res
            if (data.data.algorithms) {
                dispatch({ type: 'ALGO_LIST', payload: data.data });
            }
        })
        .catch(err =>{
            data = err.response
        });
}
export async function getRequests(dispatch,payload){
    dispatch({ type: 'REQUESTS_LOADING'});
    var data = {};
    let config = {
        headers: {'Authorization': `Bearer ${payload.token}`}
    }
        await axios.get(`${API_URL}/request/user/`+payload.userID,config).then(res => {
            data = res
            if (data.data) {
                dispatch({ type: 'REQUESTS_LIST', payload: data.data });
            }
        })
        .catch(err =>{
            if(err.response.status===401){
                alert("Your session has expired, please log in.")
                logout(dispatch)
            }
        });
}
export async function getDatas(dispatch,payload){
    dispatch({ type: 'DATAS_LOADING'});
    var data = {};
    let config = {
        headers: {'Authorization': `Bearer ${payload.token}`}
    }
        await axios.get(`${API_URL}/data/user/`+payload.userID,config).then(res => {
            data = res
            if (data.data) {
                dispatch({ type: 'DATAS_LIST', payload: data.data });
            }
        })
        .catch(err =>{
            if(err.response.status===401){
                alert("Your session has expired, please log in.")
                logout(dispatch)
            }
        });
}
export async function getAlgorithmHelp(dispatch,payload){
    var data = {};
        await axios.get(`${API_URL}/help/`+payload.algorithm).then(res => {
            data = res
            if (data.data.algorithm) {
                dispatch({ type: 'ADD_HELP', payload: data.data });
            }
        })
        .catch(err =>{
            data = err.response
        });
    return data.data
}
export async function deleteRequest(dispatch,payload){
    dispatch({ type: 'DATAS_LOADING'});
    var data = {};
    let config = {
        headers: {'Authorization': `Bearer ${payload.token}`}
    }
        await axios.delete(`${API_URL}/request/`+payload.id,config).then(res => {
            data = res
            if (data.status===204) {
                dispatch({ type: 'DELETE_REQUEST', payload: payload.id });
            }
        })
        .catch(err =>{
            if(err.response.status===401){
                alert("Your session has expired, please log in.")
                logout(dispatch)
            }
        });
}
export async function deleteData(dispatch,payload){
    dispatch({ type: 'DATAS_LOADING'});
    var data = {};
    let config = {
        headers: {'Authorization': `Bearer ${payload.token}`}
    }
        await axios.delete(`${API_URL}/data/`+payload.id,config).then(res => {
            data = res
            if (data.status===204) {
                dispatch({ type: 'DELETE_DATA', payload: payload.id });
            }
        })
        .catch(err =>{
            if(err.response.status===401){
                alert("Your session has expired, please log in.")
                logout(dispatch)
            }
        });
}

export async function newRequest(dispatch,payload) {
    dispatch({ type: 'REQUEST_LOADING' });
    let config = {}
    var logged = false
    if(payload.token){
        config = {
            headers: 
            {'Authorization': `Bearer ${payload.token}`,'Content-Type':'application/json'} 
          }
          logged=true
    }else{
        config = {
            headers: 
            {'Content-Type':'application/json'}    
          }
    }
    var data = {};
    var input = JSON.stringify({input:payload.input,parameters:[]}) 
    if(!isEmpty(payload.parameters)){
        input = JSON.stringify({input:payload.input,parameters:payload.parameters.replace(/ /g, '').split(",")})  
    } 
    await axios.post(`${API_URL}/request/`+payload.algorithm,input,config
    ).then(res => {
        if(res.status===200){
            data = res.data
            if(logged){
                dispatch({ type: 'NEW_REQUEST', payload:data.data });
            }
        }
    })
    .catch(err =>{
        if(err.response.status===401){
            alert("Your session has expired, please log in.")
        }
    });  
    dispatch({ type: 'REQUEST_LOADING_DONE' });          
return data ;
};

export async function newRequestWithData(dispatch,payload) {
    dispatch({ type: 'REQUEST_LOADING' });
    let config = {}
    config = {
        headers: 
        {'Authorization': `Bearer ${payload.token}`,'Content-Type':'application/json'} 
        }
    var input = JSON.stringify({parameters:[]}) 
    if(!isEmpty(payload.parameters)){
        input = JSON.stringify({parameters:payload.parameters.replace(/ /g, '').split(",")})  
    }  
    var data = {};
    await axios.post(`${API_URL}/request/`+payload.algorithm+"/"+payload.input,input,config
    ).then(res => {
        if(res.status===200){
            data = res.data
            dispatch({ type: 'NEW_REQUEST', payload:data.data });
        }
    })
    .catch(err =>{
        if(err.response.status===401){
            alert("Your session has expired, please log in.")
        }
    });            
    dispatch({ type: 'REQUEST_LOADING_DONE' });
return data ;
};

export async function newData(dispatch,payload) {
    let config = {}
    config = {
        headers: 
        {'Authorization': `Bearer ${payload.token}`,'Content-Type':'application/json'} 
    }
    var data = {};
    var input = JSON.stringify({data:payload.data}) 
    await axios.post(`${API_URL}/data`,input,config
    ).then(res => {
        if(res.status===200){
            data = res.data
            dispatch({ type: 'NEW_DATA', payload:data});
        }
    })
    .catch(err =>{
        if(err.response.status===401){
            alert("Your session has expired, please log in.")
        }
    });            
return data ;
};

export async function updateData(dispatch,payload){
    var data = {message:"Data update failed"};
    let config = {
        headers: {'Authorization': `Bearer ${payload.token}`}
    }
    var input = JSON.stringify(payload.updateBody) 
        await axios.put(`${API_URL}/data/`+payload.dataID,input,config).then(res => {
            if (res.status===200) {
                data = {message:"OK"}
            }
        })
        .catch(err =>{
            if(err.response.status===401){
                alert("Your session has expired, please log in.")
                logout(dispatch)
            }else if(err.response.status===400){
                data = err.response.data
            }else{
                data = {message:"Data update failed"}
            }
        });
    return data
}

export async function logout(dispatch) {
    dispatch({ type: 'LOGOUT' });
};