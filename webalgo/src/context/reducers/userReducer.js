import Cookies from 'js-cookie'
 let user = Cookies.get('user')
     ? JSON.parse(Cookies.get('user')).userID
     : "";
     let token = Cookies.get('user')
     ? JSON.parse(Cookies.get('user')).token
     : "";
     let username = Cookies.get('user')
     ? JSON.parse(Cookies.get('user')).username
     : "";

const initialState = {
    userID: "" || user,
    token: "" || token,
    username: username || "Guest",
    email: "",
    loadingRequest: false,
    loadingRequests: false,
    loadingDatas: false,
    errorMessage: null,
    requests:[],
    datas:[],
    algorithms:[],
    helps:[]
};

const userReducer = (state = initialState,action) =>{
   const {type,payload} = action;
   switch (type) {
       case "LOGIN_REQUEST":
           return {
               ...state, 
               loading:true
           };
       case "LOGIN_OK":
       return {
           ...state,
           userID: payload.userID,
           token: payload.token,
           username: payload.username,
           loading:false
       };
       case "LOGOUT":
           Cookies.remove('user')
       // localStorage.removeItem('currentUser');
       return {
           ...state,
            userID: "",
            token: "",
            username: "Guest",
            email: "",
            loadingRequest: false,
            loadingRequests: false,
            loadingDatas: false,
            errorMessage: null,
            requests:[],
            datas:[],
            helps:[]
       };
       case "LOGIN_ERROR":
       return {
           ...state,
           loading: false,
           errorMessage: "Login error"
       };
       case "USER_INFO":
           return {
            ...state,
            username: payload.username,
            email: payload.email
           }
        case "ALGO_LIST":
            return {
            ...state,
            algorithms: payload.algorithms
            };
        case "REQUESTS_LIST":
            return {
            ...state,
            loadingRequests:false,
            requests: payload
            }
        case "REQUEST_LOADING":
            return {
            ...state,
            loadingRequest: true
        }
        case "REQUEST_LOADING_DONE":
            return {
            ...state,
            loadingRequest: false
        }
        case "REQUESTS_LOADING":
            return {
            ...state,
            loadingRequests: true
        }
        case "DATAS_LIST":
            return {
            ...state,
            loadingDatas:false,
            datas: payload
            }
        case "DATAS_LOADING":
            return {
            ...state,
            loadingDatas: true
        }
        case "ADD_HELP":
            return {
                ...state,
                helps: [...state.helps, payload]
            }
            case "DELETE_REQUEST":
            return {
                ...state,
                requests: [...state.requests.filter(request => request.id !==payload)]
            }
            case "DELETE_DATA":
            return {
                ...state,
                datas: [...state.datas.filter(data => data.id !==payload)]
            }
        case "NEW_REQUEST":
            return{
                ...state,
                requests: [...state.requests,payload]
            }
        case "NEW_DATA":
            return{
                ...state,
                datas: [...state.datas,payload]
            }
        case "ADD_EMAIL":
            return{
                ...state,
                email:payload.email            }
       default:
           return state;
   }
}

export default userReducer;