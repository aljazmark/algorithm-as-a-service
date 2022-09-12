import React, { useEffect } from 'react';
import { CircularProgress,FormControlLabel,TextField,Button,Checkbox,MenuItem,Select,InputLabel,makeStyles,FormControl } from '@material-ui/core';
import './css/newRequest.css'
import {newRequest,newRequestWithData} from '../context'
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
    algorithm:"",
    input:"",
    parameters:[],
    withData:false
}


export const NewRequest = () => {
    const classes = useStyles();
    const [formValues,setFormValues] = React.useState(initialFormValues);
    const algorithms = useSelector((state)=> state.user.algorithms);
    const [ result,setResult] = React.useState({}); 
    const [show,setShow] = React.useState(false);
    const [displayed,setDisplayed] = React.useState("");
    const token = useSelector((state)=> state.user.token);
    const loadingRequest = useSelector((state)=> state.user.loadingRequest);
    const dispatchR = useDispatch()
    var userLogged = (token) ? true : false;
    if(token==="" && formValues.withData){
        setFormValues({...formValues,withData:false})
    }
    const handleChange= e =>{
        const { name,value} = e.target
        setFormValues({
            ...formValues,
            [name]:value
        })
    }
    const handleRequest = () => {
        setShow(false)
        var tmpProps = {token:token,algorithm:formValues.algorithm,input:formValues.input,parameters:formValues.parameters}
        if(formValues.withData){  
            newRequestWithData(dispatchR,tmpProps).then(res=> {setResult(res)})
        }else{
            newRequest(dispatchR,tmpProps).then(res=> {setResult(res)})
        }
        
    }
    useEffect(() => {
        if(!isEmpty(result)){
            setShow(true)
        }
    }, [result]);

    const showLoading = () => {
        return (
            <div>
                <CircularProgress size={60}/>
            </div>
        )
    }
    useEffect(() => {
        if(loadingRequest){
            setDisplayed("loading")
        }else if(!loadingRequest&&show){
            setDisplayed("result")
        }
    }, [loadingRequest,show]);
    const showResult = () => {
        return (
        <div>
            <h3>Request ID</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={result.id}
                InputProps={{
                    readOnly: true,
                  }}
            ></TextField>
            <h3>Algorithm</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={result.algorithm}
                InputProps={{
                    readOnly: true,
                  }}
            ></TextField>
            <h3>Input</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={result.input}
                InputProps={{
                    readOnly: true,
                  }}
            ></TextField>
            <h3>Parameters</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={JSON.stringify(result.parameters)}
                InputProps={{
                    readOnly: true,
                  }}
            ></TextField>
            <h3>Output</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={result.output}
                InputProps={{
                    readOnly: true,
                  }}
            ></TextField>
            <h3>Execution time</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={result.executionTime}
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
                <h1 data-testid="requests-title">New request</h1>
                <form>
                    <FormControl variant="outlined">
                        <InputLabel>Select an algorithm</InputLabel>
                        <Select
                            variant="outlined"
                            label="Select an algorithm"
                            name="algorithm"
                            className={classes.resultText}
                            value={formValues.algorithm}
                            onChange= {handleChange}>
                            {
                                algorithms.map(
                                    item => (<MenuItem key={item} value={item}>{item}</MenuItem>)

                                )
                            }
                        </Select>
                        <TextField
                            variant="outlined"
                            label="Enter your input"
                            name="input"
                            className={classes.inputField}
                            value={formValues.input}
                            onChange={handleChange}>
                        </TextField>
                        {userLogged ?
                        <FormControlLabel
                        control={<Checkbox checked={formValues.withData} onChange={(e) => {handleChange({target:{name:e.target.name,value:e.target.checked}})}} name="withData" />}
                        label="Use data ID as input"
                        />
                        :
                        null
                        }
                        
                        <TextField
                            variant="outlined"
                            label="Enter your parameters"
                            name="parameters"
                            className={classes.inputField}
                            value={formValues.parameters}
                            helperText="Format: parameter 1, parameter 2,..."
                            onChange={handleChange}>
                        </TextField>
                        <div style={{padding:"30px"}}>
                        <Button disabled={loadingRequest} variant="outlined" size="large" onClick={()=>handleRequest()}>Submit</Button>
                        </div>
                    </FormControl>    
                </form>
            </div>
            <div className="divTwo" align="center">
            <h1 data-testid="requests-title">Result</h1>
            {{
                "": null,
                "loading": showLoading(),
                "result": showResult(),
                default: (
                null
                )
            }[displayed]}
            </div>
        </div>
    )
}