
import React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import {makeStyles,FormControl} from '@material-ui/core'
import { useDispatch } from "react-redux";
import store from '../context/store'
import {updateData,getDatas} from '../context'
  
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

export const DataDialog = (props)=> {
  var userID = store.getState().user.userID
  var token = store.getState().user.token
  const [formError,setFormError]= React.useState("")
  const [open, setOpen] = React.useState(false);
  const classes = useStyles();
  var dataID = props.dataID
  var data = props.data
  var allowed = ""
  if(props.allowed.length > 0){
    props.allowed.forEach(function(item,index){
        if(index!==props.allowed.length-1){
            allowed=allowed+item+","
        }else if(index===props.allowed.length-1){
            allowed=allowed+item
        }
    });
  }
  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setFormError("")
  };

  const dispatchR = useDispatch();
  
  const handleEdit = ()=>{
      updateData(dispatchR,{token:token,dataID:dataID,updateBody:{data:data,allowed:allowed.replace(/ /g, '').split(",")}}).then(res=> {
          if(res.message!=="OK"){
            setFormError(res.message)
          }else{
            setOpen(false);
            setFormError("")
            getDatas(dispatchR,{userID:userID,token:token})
          }
      })
  }
  const DataForm = () => {
      return(
        <div align="center">
            <h1 id="form-dialog-title"   align="center" >Data</h1>
            <FormControl variant="outlined">
            <DialogContent>
            <TextField
            className={classes.resultText}
            label="Data ID"
            variant="outlined"
            value={dataID}
            InputProps={{
                readOnly: true,
            }}
            ></TextField>    
            <TextField
                variant="outlined"
                label="Data"
                name="data"
                multiline
                className={classes.inputField}
                defaultValue={data}
                onChange={(e)=>{data=e.target.value}}
            />
            <TextField
                variant="outlined"
                label="Allowed"
                name="allowed"
                multiline
                className={classes.inputField}
                defaultValue={allowed}
                onChange={(e)=>{allowed=e.target.value}}
            />
            <h4>{formError}</h4>
            </DialogContent>
            </FormControl>
            <DialogActions>
            <Button onClick={handleClose} color="primary">
                Cancel
            </Button>
            <Button onClick={handleEdit} color="primary">
                Confirm
            </Button>
            </DialogActions>
        </div>
      )
  }
  return (
    <div>
      <Button color="inherit" onClick={handleClickOpen}>
        Edit
      </Button> 
      
      <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
        {<DataForm />}
      </Dialog>
    </div>
  );
}