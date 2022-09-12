import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import { withStyles} from '@material-ui/core/styles';
import { login, register,logout } from '../context';
import { useDispatch,useSelector } from "react-redux";
const SwitchButton = withStyles({
    root: {
      fontSize: '1rem',
      background: 'white',
      borderRadius: 3,
      border: 0,
      color: 'primary',
      margin: '30px 15px 30px 15px',
    },
  })(Button);


export default function LoginDialog() {
  const state = useSelector((state)=> state.user);
  const [logErr,setLogErr] = React.useState("");
  const [open, setOpen] = React.useState(false);
  const [form, setForm] = React.useState(true);
  
  let password = "";
  let username = "";
  let password2 = "";
  let email = "";
  const handleClickOpen = () => {
    setForm(true)
    password="";
    username="";
    password2 = "";
    email = "";
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    password="";
    username="";
    password2 = "";
    email = "";
    setForm(true)
  };

  const dispatchR = useDispatch();
  
  const handleLogin = async (e) =>{
    e.preventDefault();
    var response
    let payload = {username,password}
    try {
      setLogErr("Logging in, please wait")
      response = await login(dispatchR,payload);
      setLogErr("")
    }catch(error){
    }
    if(response.status!==200){
      setLogErr(response.data.message)
    }else{
      handleClose();
    }
  }
  const handleRegister = async (e) =>{
    e.preventDefault();
    if(email===""){
      setLogErr("Please enter an email")
      return
    }else if(username===""){
      setLogErr("Please enter a username")
      return
    }else if(password===""){
      setLogErr("Please enter a password")
      return
    }else if(password2===""){
      setLogErr("Please repeat the password")
      return
    }else if(password!==password2){
      setLogErr("Passwords do not match")
      return
    }
    var response
    let payload = {email,username,password}
    try {
      setLogErr("Registering, please wait")
      response = await register(dispatchR,payload);
      setLogErr("")
    }catch(error){
    }
    if(response.status!==200){
      setLogErr(response.data.message)
    }else{
      handleClose();
    }
  }
  const handleLogout = async ()=>{
    logout(dispatchR);
  }
  const LoginForm = () => {
      return(
        <div>
            <h1 id="form-dialog-title"   align="center" >Login</h1>
            <DialogContent>
            <DialogContentText>
            </DialogContentText>
            <TextField
                autoFocus
                margin="dense"
                id="username"
                label="Username"
                type="text"
                required={true}
                fullWidth
                onChange={(e)=>{username=e.target.value}}
            />
            <TextField
                margin="dense"
                id="password"
                label="Password"
                type="password"
                required={true}
                fullWidth
                onChange={(e)=>{password=e.target.value}}
            />
            <h4>{logErr}</h4>
            </DialogContent>
            <div align="center" >
                <SwitchButton onClick={() => {setForm(false)}} color="primary" >
                    Register
                </SwitchButton>
                <h3>{}</h3>
            </div>
            <DialogActions>
            <Button onClick={handleClose} color="primary">
                Cancel
            </Button>
            <Button onClick={handleLogin} color="primary">
                Login
            </Button>
            </DialogActions>
        </div>
      )
  }
  const RegisterForm = () => {
    return(
      <div>
          <h1 id="form-dialog-title"  align="center">Register</h1>
          <DialogContent>
          <DialogContentText>
          </DialogContentText>
          <TextField
              autoFocus
              margin="dense"
              id="email"
              label="Email"
              type="email"
              required={true}
              fullWidth
              onChange={(e)=>{email=e.target.value}}
          />
          <TextField
                margin="dense"
                id="username"
                label="Username"
                type="text"
                required={true}
                fullWidth
                onChange={(e)=>{username=e.target.value}}
            />
          <TextField
              margin="dense"
              id="password"
              label="Password"
              type="password"
              required={true}
              fullWidth
              onChange={(e)=>{password=e.target.value}}
          />
          <TextField
              margin="dense"
              id="password2"
              label="Confirm password"
              type="password"
              required={true}
              fullWidth
              onChange={(e)=>{password2=e.target.value}}
          />
          <h4>{logErr}</h4>
          </DialogContent>
          <div align="center"><SwitchButton onClick={() => {setForm(true)}} color="primary" >
                Login
            </SwitchButton></div>
          
          <DialogActions>
          <Button onClick={handleClose} color="primary">
              Cancel
          </Button>
          <Button onClick={handleRegister} color="primary">
              Register
          </Button>
          </DialogActions>
      </div>
    )
}
  return (
    <div>
      {state.userID ? <Button color="inherit" onClick={handleLogout}>
        Logout
      </Button> 
      :
      <Button color="inherit" onClick={handleClickOpen}>
        Login
      </Button>}
      
      <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
        {form ? (<LoginForm />) : (<RegisterForm />)}
      </Dialog>
    </div>
  );
}