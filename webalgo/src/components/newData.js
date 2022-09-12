import React, { useEffect } from 'react';
import {TextField,Button,makeStyles,FormControl } from '@material-ui/core';
import './css/newRequest.css'
import {newData} from '../context'
import {useSelector,useDispatch} from 'react-redux'
import {isEmpty} from "lodash"
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
    data:""
}


export const NewData = () => {
    const classes = useStyles();
    const [formValues,setFormValues] = React.useState(initialFormValues);
    const [ result,setResult] = React.useState({}); 
    const [show,setShow] = React.useState(false);
    const token = useSelector((state)=> state.user.token);
    const dispatchR = useDispatch()
    const handleChange= e =>{
        const { name,value} = e.target
        setFormValues({
            ...formValues,
            [name]:value
        })
    }
    const handleRequest = () => {
        var tmpProps = {token:token,data:formValues.data}
        newData(dispatchR,tmpProps).then(res=> {setResult(res)})
        
    }
    useEffect(() => {
        if(!isEmpty(result)){
            setShow(true)
        }
    }, [result]);

    

    const showResult = () => {
        return (
        <div>
            <h3>Data ID</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={result.id}
                InputProps={{
                    readOnly: true,
                  }}
            ></TextField>
            <h3>Data</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={result.data}
                InputProps={{
                    readOnly: true,
                  }}
            ></TextField>
        </div>
                )
    }
    return (
        <div className="flexDir" >
            <div className="divOne" data-testit="requests" align="center">
                <h1 data-testid="requests-title">New data</h1>
                <form>
                    <FormControl variant="outlined">
                        <TextField
                            variant="outlined"
                            label="Enter your data"
                            name="data"
                            className={classes.inputField}
                            value={formValues.data}
                            onChange={handleChange}>
                        </TextField>
                        <div style={{padding:"30px"}}>
                        <Button variant="outlined" size="large" onClick={()=>handleRequest()}>Submit</Button>
                        </div>
                    </FormControl>    
                </form>
            </div>
            <div className="divTwo" align="center">
            <h1 data-testid="requests-title">Result</h1>
            {show ? showResult(): null}
            </div>
        </div>
    )
}