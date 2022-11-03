import React, { useState, useEffect } from 'react';
import { Col, Row, Card, Button, InputGroup, FormControl } from 'react-bootstrap';
import { CircularProgressbar, buildStyles } from "react-circular-progressbar";
import { BsCpu } from 'react-icons/bs';
import { IoHardwareChipOutline } from 'react-icons/io5';
import "react-circular-progressbar/dist/styles.css";
import useFetch from 'use-http'
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { API_ENDPOINT, WS_ENDPOINT } from "../Global"
import { GetAccessToken } from "../AAD"

import { Fade } from "react-awesome-reveal";
import { Section, AnimatedGradientH1 } from './Styled';

const Stats = () => {
  const [realTimeStats, setRealTimeStats] = useState();
  const [stats, setStats] = useState();
  const [rtt, setRTT] = useState(0);
  const [cpuLoad, setCpuLoad] = useState(0)
  const [memLoad, setMemLoad] = useState(0)


  const { response, get, post } = useFetch(API_ENDPOINT, {
    cachePolicy: "no-cache",
    retries: 2,
    retryDelay: ({ attempt, error, response }) => {
      return Math.min(attempt > 1 ? 2 * attempt * 1000 : 1000, 30 * 1000)
    },
    interceptors: {
      request: async ({ options, url, path, route }) => {
        let accessToken = await GetAccessToken()
        options.headers.Authorization = "Bearer " + accessToken
        return options
      },
      response: async ({ response }) => {
        return response
      }
    },
  }
  )

  const {
    lastJsonMessage,
    readyState,
  } = useWebSocket(WS_ENDPOINT,
    {
      shouldReconnect: (closeEvent) => true,
    });

  useEffect(() => {
    const getStats = async () => {
      const stats = await get('/v1/system/stats')
      if (response.ok) {
        setStats(stats)
      }
    }
    getStats()
  }, [get, response.ok]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      setRealTimeStats(lastJsonMessage);
      console.log(lastJsonMessage)
    }
  }, [lastJsonMessage, setRealTimeStats]);

  useEffect(() => {
    const id = setInterval(async () => {
      let accessToken = await GetAccessToken()
      var start = Date.now()
      fetch(API_ENDPOINT + "/v1/ping", {
        headers: { Authorization: "Bearer " + accessToken }
      }).then(response => {
        var latency = Date.now() - start
        setRTT(latency)
      })
    }, 30000)
    return () => clearInterval(id);
  }, []);

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Connecting',
    [ReadyState.OPEN]: 'Connected',
    [ReadyState.CLOSING]: 'Closing',
    [ReadyState.CLOSED]: 'Closed',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  async function updateCpuLoad() {
    await post('/v1/system/cpu', { "percent": parseInt(cpuLoad) })
    if (response.ok) {
      console.log("cpu load posted")
    }
  }

  async function updateMemoryLoad() {
    await post('/v1/system/memory', { "size": parseInt(memLoad) })
    if (response.ok) {
      console.log("memory load posted")
    }
  }

  function formatUptime(seconds) {
    const days = Math.floor(seconds / 3600 / 24)
    seconds %= (3600 * 24);
    const hours = Math.floor(seconds / 3600);
    seconds %= 3600;
    const minutes = Math.floor(seconds / 60);
    seconds = seconds % 60;
    return days + " Days, " + hours.toString().padStart(2, "0") + ":" + minutes.toString().padStart(2, "0") + ":" + seconds.toString().padStart(2, "0")
  }

  return (
    <Section color="#f7f5eb">
      <AnimatedGradientH1 className="mb-4">System Statistics</AnimatedGradientH1>
      <Row className='row-cols-1 row-cols-md-2 row-cols-xl-3 gy-4'>
        <Col>
          <Card className='shadow h-100'>
            <Card.Header className="text-center text-primary"><h6>CPU</h6></Card.Header>
            <Card.Body>
              <Fade direction="down" cascade={true}>
                <Row>
                  <Col className='col-4 align-self-center'>
                    <h6>{stats?.cpu_cores ? stats?.cpu_cores : "N/A"}</h6><Card.Text>Cores</Card.Text>
                    <h6>{stats?.cpu_cache_size ? stats?.cpu_cache_size + " KB" : "N/A"}</h6><Card.Text>Cache</Card.Text>
                  </Col>
                  <Col className='col-4'>
                    <h6 className="text-center">Utilization</h6>
                    <CircularProgressbar
                      value={(realTimeStats?.cpu_percent ? Math.round(realTimeStats?.cpu_percent * 100) / 100 : 0)}
                      text={(realTimeStats?.cpu_percent ? Math.round(realTimeStats?.cpu_percent * 100) / 100 : 0) + "%"}
                      strokeWidth={3}
                      circleRatio={0.75}
                      styles={buildStyles({
                        rotation: 1 / 2 + 1 / 8,
                        textSize: "1.5em",
                        textColor: "#4da6e7",
                        pathColor: "#4da6e7",
                        trailColor: "#eee"
                      })}
                    >
                    </CircularProgressbar>
                  </Col>
                  <Col className='col-4 align-self-center'>
                    <BsCpu size={28} />
                    <Card.Text className="small">{stats?.cpu_model_name ? stats?.cpu_model_name : "N/A"}</Card.Text>
                  </Col>
                </Row>
                <Row className="mt-4">
                  <Col >
                    <InputGroup>
                      <InputGroup.Text>+</InputGroup.Text>
                      <FormControl value={cpuLoad} onInput={e => setCpuLoad(e.target.value)} type="number" min="0" max="100" placeholder={cpuLoad} />
                      <InputGroup.Text>%CPU</InputGroup.Text>
                      <Button className='rounded ms-4 ps-4 pe-4' onClick={updateCpuLoad}>
                        Go
                      </Button>
                    </InputGroup>
                  </Col>
                </Row>
              </Fade>
            </Card.Body>
          </Card>
        </Col>
        <Col>
          <Card className='shadow h-100'>
            <Card.Header className="text-center text-primary"><h6>MEMORY</h6></Card.Header>
            <Card.Body>
              <Fade direction="down" cascade={true}>
                <Row>
                  <Col className='col-4 align-self-center'>
                    <h6>{realTimeStats?.mem_used ? Math.round(realTimeStats?.mem_used / 1024 / 1024 * 100) / 100 + " MB" : "N/A"}</h6><Card.Text>Used</Card.Text>
                    <h6>{stats?.mem_total ? Math.round(stats?.mem_total / 1024 / 1024 * 100) / 100 + " MB" : "N/A"}</h6><Card.Text>Total</Card.Text>
                  </Col>
                  <Col className='col-4'>
                    <h6 className="text-center">Usage</h6>
                    <CircularProgressbar
                      value={realTimeStats?.mem_percent ? Math.round(realTimeStats?.mem_percent * 100) / 100 : 0}
                      text={(realTimeStats?.mem_percent ? Math.round(realTimeStats?.mem_percent * 100) / 100 : 0) + "%"}
                      strokeWidth={3}
                      circleRatio={0.75}
                      styles={buildStyles({
                        rotation: 1 / 2 + 1 / 8,
                        textSize: "1.5em",
                        textColor: "#726ae3",
                        pathColor: "#726ae3",
                        trailColor: "#eee"
                      })}
                    >
                    </CircularProgressbar>
                  </Col>
                  <Col className='col-4 align-self-center'>
                    <IoHardwareChipOutline size={28} />
                  </Col>
                </Row>
                <Row className="mt-4">
                  <Col >
                    <InputGroup>
                      <InputGroup.Text>+</InputGroup.Text>
                      <FormControl value={memLoad} onInput={e => setMemLoad(e.target.value)} type="number" min="0" max="16777216" placeholder={memLoad} />
                      <InputGroup.Text>MB</InputGroup.Text>
                      <Button onClick={updateMemoryLoad} className='rounded ms-4 ps-4 pe-4'>
                        Go
                      </Button>
                    </InputGroup>
                  </Col>
                </Row>
              </Fade>
            </Card.Body>
          </Card>
        </Col>
        <Col>
          <Card className='shadow h-100'>
            <Card.Header className="text-center text-primary"><h6>NETWORK</h6></Card.Header>
            <Card.Body>
              <Fade direction="down" cascade={true}>
                <Row>
                  <Col className='col-4 align-self-center'>
                    <h6>{connectionStatus}</h6><Card.Text>Web socket</Card.Text>
                    <h6>{rtt} ms</h6><Card.Text>RTT</Card.Text>
                  </Col>
                  <Col className='col-4'>
                    <h6 className="text-center">Latency</h6>
                    <CircularProgressbar
                      value={3}
                      text={rtt + "ms"}
                      maxValue={500}
                      circleRatio={0.75}
                      strokeWidth={3}
                      styles={buildStyles({
                        textSize: "1.5em",
                        textColor: "#f58b56",
                        pathColor: "#f58b56",
                        rotation: 1 / 2 + 1 / 8,
                        trailColor: "#eee"
                      })}
                    />
                  </Col>
                  <Col className='col-4 align-self-center'>
                    <h6>{stats?.local_ip ? stats?.local_ip : "N/A"}</h6><Card.Text>Internal IP</Card.Text>
                    <h6>{stats?.external_ip ? stats?.external_ip : "N/A"}</h6><Card.Text>External IP</Card.Text>
                  </Col>
                </Row>
              </Fade>
            </Card.Body>
          </Card>
        </Col>
      </Row>
      <Row className="mt-4">
        <Col>
          <Card className='shadow'>
            <Card.Header className="text-center text-primary"><h6>SYSTEM{stats?.environment ? "("+stats?.environment+")" : ""}</h6></Card.Header>
            <Card.Body>
              <Row>
                <Col className='col-12 col-md-4 align-self-center'>
                  <h6>Platform: {stats?.platform ? stats?.platform : "N/A"}-{stats?.platform_version ? stats?.platform_version : "N/A"}</h6>
                  <h6>OS: {stats?.os ? stats?.os : "N/A"}</h6>
                  <h6>Hostname: {stats?.hostname ? stats?.hostname : "N/A"}</h6>
                  <h6>Uptime: {realTimeStats?.uptime ? formatUptime(realTimeStats?.uptime) : "N/A"}</h6>
                </Col>
                <Col className='col-12 col-md-4 align-self-center'>
                  <h6>Go Version: {stats?.go_version ? stats?.go_version : "N/A"}</h6>
                  <h6>Arch: {stats?.go_arch ? stats?.go_arch : "N/A"}</h6>
                  <h6>Docker: {stats?.is_docker ? "Yes" : "No"}</h6>
                  <h6>Kubernetes: {stats?.is_kubernetes ? "Yes" : "No"}</h6>
                </Col>
                <Col className='col-12 col-md-4 align-self-center'>
                  <h6>Location: {stats?.location ? stats?.location : "N/A"}</h6>
                  <h6>Zone: {stats?.zone ? stats?.zone : "N/A"}</h6>
                  <h6>Name: {stats?.name ? stats?.name : "N/A"}</h6>
                  <h6>Size: {stats?.size ? stats?.size : "N/A"}</h6>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Section>
  );
};
export default Stats;
