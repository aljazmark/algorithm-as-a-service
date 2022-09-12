import React, { useEffect } from 'react';
import {TextField,Button,makeStyles,FormControl } from '@material-ui/core';
import './css/newRequest.css'
import {useDispatch} from 'react-redux'
import store from '../context/store'
import {ConfirmDialog} from './confirmDialog'
import {updateUser} from '../context'
import {useSelector} from "react-redux";
const useStyles = makeStyles((theme) => ({
    formControl: {
      margin: theme.spacing(0),
      minWidth: 210,
      minHeight: 150,
    },
    selectEmpty: {
      marginTop: theme.spacing(0),
    },
    resultText:{
        minWidth:300
    },
    inputField:{
        minWidth:300,
        marginTop:theme.spacing(3)
    }
  }));
const initialFormValues ={
    id:"",
    username:"",
    email:"",
    password:"",
    password2:""
}


export const Settings = (props) => {
    const classes = useStyles();
    const [formValues,setFormValues] = React.useState(initialFormValues);
    const email = useSelector((state)=> state.user.email);
    var username = store.getState().user.username
    const [openDialog,setOpenDialog]= React.useState(false)
    const [confirmChange,setConfirmChange]= React.useState(false)
    const [formError,setFormError]= React.useState("")
    const dispatchR = useDispatch()
    const handleChange= e =>{
        const { name,value} = e.target
        setFormValues({
            ...formValues,
            [name]:value
        })
    }
    useEffect(() => {
        //initialize values
        setFormValues({
            ...formValues,
            id:props.userID,
            username:username,
            email:email
        })
    }, [email]); // eslint-disable-line react-hooks/exhaustive-deps
    
    const handleRequest = () => {
        setOpenDialog(true)
    }
    useEffect(() => {
        //setting values
        if(confirmChange){
            setConfirmChange(false)
            setOpenDialog(false)
            if(formValues.password!=="" && formValues.password===formValues.password2){
                updateUser(dispatchR,{token:props.token,userID:props.userID,updateBody:{username:formValues.username,email:formValues.email,password:formValues.password}}).then(res=> {setFormError(res.message)})
            }else if(formValues.password===""){
                updateUser(dispatchR,{token:props.token,userID:props.userID,updateBody:{username:formValues.username,email:formValues.email}}).then(res=> {setFormError(res.message)})
            }else{
                setFormError("Passwords must match")
            }
        }
    }, [confirmChange]); // eslint-disable-line react-hooks/exhaustive-deps
    return (
        <div className="flexDir" >
            <div className="divOne" data-testit="requests" align="center">
                <h1 data-testid="requests-title">User settings</h1>
                <form>
                    
                    <FormControl variant="outlined">
                        <TextField
                            className={classes.resultText}
                            label="User ID"
                            variant="outlined"
                            value={formValues.id}
                            InputProps={{
                                readOnly: true,
                            }}
                        ></TextField>
                        <TextField
                            variant="outlined"
                            label="Username"
                            name="username"
                            className={classes.inputField}
                            value={formValues.username}
                            onChange={handleChange}>
                        </TextField>
                        <TextField
                            variant="outlined"
                            label="Email"
                            name="email"
                            className={classes.inputField}
                            value={formValues.email}
                            onChange={handleChange}>
                        </TextField>
                        <TextField
                            variant="outlined"
                            label="New password"
                            name="password"
                            className={classes.inputField}
                            value={formValues.password}
                            type="password"
                            onChange={handleChange}>
                        </TextField>
                        <TextField
                            variant="outlined"
                            label="Repeat new password"
                            name="password2"
                            className={classes.inputField}
                            value={formValues.password2}
                            type="password"
                            onChange={handleChange}>
                        </TextField>
                        <h4>{formError}</h4>
                        <div style={{padding:"30px"}}>
                        <Button variant="outlined" size="large" onClick={()=>handleRequest()}>Submit</Button>
                        <ConfirmDialog 
                        title="Please confirm to proceed"
                        content="You are about to change your user settings. Doing so will log you out of the service."
                        confirm={openDialog}
                        setConfirm={setOpenDialog}
                        confirmChange={setConfirmChange}
                        />
                        </div>
                    </FormControl>    
                </form>
            </div>
        </div>
    )
}