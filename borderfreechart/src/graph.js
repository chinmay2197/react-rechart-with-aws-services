//this line imports react functionality
import React from 'react';
import { useEffect, useState } from 'react';
import {
    LineChart, Line, XAxis, YAxis, Tooltip, Legend, CartesianGrid
  } from 'recharts';
  import Moment from 'moment';

const datamap = {}
function CustomTooltip(props) {
  var Interest = '';
  var Time = '';

  if (props.label){
    Interest = props.data[props.label];
    Time = props.label;
  }

  return (
    <div style={{
      border:"1px solid rgb(204, 204, 204)",
      height:"50px",
      width:"100px",
      backgroundColor: "rgb(255, 255, 255, 255)"
    }}>
      <div style={{
        position: 'absolute',
        left: '20%',
        top: '50%',
        transform: 'translate(-20%, -50%)'
      }}>
      <div> {Moment(Time).format('MMM YYYY')} </div>
      <div style={{ color: '#8884d8' }}> Score: {Interest}</div>
      </div>
    </div>
  );
}

function formatXAxis(tickItem){
  return Moment(tickItem).format('MMM 1,YYYY')
}

export default function App() {
  const [isLoaded, setIsLoaded] = useState(true);
  const [items, setItems] = useState([]);
  const data = [];

  useEffect(() => {
      const api = 'https://34cnwmw3ic.execute-api.ap-south-1.amazonaws.com/dev'
    fetch(
      api, {
        mode: 'cors',
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((res) => res.json())
      .then((result) => {
          for (var instance in result['items']) {
            var mydata = (result['items'][instance]);
            mydata.date = Moment(result['items'][instance]['Time']).valueOf();
            datamap[mydata.date] = result['items'][instance]['Interest']
            data.push(mydata);
          }
          setItems(data);
        },
        (error) => {
          setIsLoaded(false);
          console.log(error)
        },
      );
  },[]);
  if (isLoaded){
    return (
      <div className="chart"  style={{
          position: 'absolute', left: '50%', top: '50%',
          transform: 'translate(-50%, -50%)',
          backgroundColor: 'white',
          boxShadow: "rgba(0, 0, 0, 0.12) 0px 10px 10px 0px"
        }}>
        <div className="header">
          <h3 style={{
            paddingLeft: "500px",
            borderBottom:"1px solid rgb(204, 204, 204)",
            paddingBottom:"5px"
          }}>Interest Over Time</h3>
        </div>
        <br/>
        <br/>
        <LineChart
          width={1200}
          height={400}
          data={items}
          margin={{top: 2, right: 50, left: 20, bottom: 5}}
        >
          <CartesianGrid strokeDasharray="3 3" vertical={false} />
          <XAxis dataKey="date" axisLine={false} tickFormatter={formatXAxis} scale="time" type="number" interval={items.length/5} domain={["auto","auto"]} padding={{left:5, right:5, bottom:25}}/>
          <YAxis axisLine={false} type="number" domain={[0,100]} tickLine={false}/>
          <Tooltip content={<CustomTooltip data={datamap} />} />
          <Legend />
          <Line type="monotone" dataKey="Interest" stroke="#8884d8" strokeWidth={2} dot={false}/>
        </LineChart>
        <br/>
      </div>
    );
  }
  return (
    <div className="chart"  style={{
        position: 'absolute', left: '50%', top: '50%',
        transform: 'translate(-50%, -50%)',
        backgroundColor: 'white',
        boxShadow: "rgba(0, 0, 0, 0.12) 0px 10px 10px 0px"
      }}>
      <div className="header">
        <h3 style={{
          paddingLeft: "500px",
          borderBottom:"1px solid rgb(204, 204, 204)",
          paddingBottom:"5px"
        }}>Interest Over Time</h3>
      </div>
      <br/>
      <br/>
      <LineChart
        width={1200}
        height={400}
        data={items}
        margin={{top: 2, right: 50, left: 20, bottom: 5}}
      >
        <CartesianGrid strokeDasharray="3 3" vertical={false} />
        <XAxis dataKey="Time" axisLine={false} tickFormatter={formatXAxis} interval={5} padding={{left:5, right:5, bottom:25}}/>
        <YAxis axisLine={false} type="number" domain={[0,100]} tickLine={false}/>
        <Tooltip content={<CustomTooltip data={datamap} />} />
        <Legend />
      </LineChart>
      <br/>
    </div>
  )
}