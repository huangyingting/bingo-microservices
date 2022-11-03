import React from 'react';
import { HeroHeader, Paragraph } from './Styled';

const Home = () => {
  return (
    <>
      <HeroHeader color="#e1f4fd" title="Cloud Native Architecture Design" 
        subTitle="Empowering people in digital transformation"
        content="Reinvent business in the cloud.Cloud transformation processes help organizations move their information to cloud environments. This process can be scaled according to business needs."
        image="/images/company.svg"
      />
      <Paragraph color="#fff" title="High redundancy"
        content="All data centers we use have top level of redundancy of critical components. Electricity outages are prevented by multiple power feeds, on-site power generators and enterprise-class UPS technology. Uninterrupted network connectivity is guaranteed by the simultaneous usage of multiple major carriers. Redundant and geographically distributed backups across countries and continents are possible thanks to the availability of multiple locations of the Cloud’s facilities."
        image="/images/ha.svg"
      />
      <Paragraph color="#f5f6fa" title="Automatic Scaling on Cloud" reverse
        content="Autoscaling provides users with an automated approach to increase or decrease the compute, memory or networking resources they have allocated, as traffic spikes and use patterns demand."
        image="/images/autoscale.svg"
      />
      <Paragraph color="#fff" title="Latest technologies integrated" 
        content="We are consistently among the first hosting companies to provide their users access to the latest speed technologies. Our customers do not have to wait in order to take advantage of the newest PHP versions or the most innovative protocols and compression algorithms like Brotli, HTTP/2, TLS 1.3, OCSP Stapling and QUIC."
        image="/images/technologies.svg"
      />
    </>
  );
};
export default Home;
