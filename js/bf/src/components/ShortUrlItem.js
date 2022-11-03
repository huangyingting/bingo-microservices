import React from 'react'
import { RiFileCopyLine } from 'react-icons/ri';
import { FaQrcode } from 'react-icons/fa';
import { BiBarChart } from "react-icons/bi";
import { RiDeleteBinLine, RiTwitterLine, RiShareBoxLine } from "react-icons/ri";
import { FiEdit } from "react-icons/fi";
import { AiOutlineFacebook, AiOutlineWhatsApp, AiOutlineLinkedin, AiOutlineMail } from "react-icons/ai";
import { FacebookShareButton, LinkedinShareButton, WhatsappShareButton, TwitterShareButton, EmailShareButton } from "react-share";
import { Button, Tooltip, OverlayTrigger, Dropdown, Badge, Row, Col } from 'react-bootstrap';
import { GradientText } from './Styled';

function ShortUrlItem(props) {

  const domainName = (url) => {
    let domain
    try {
      domain = new URL(url)
    } catch (_) {
      return ""
    }
    return domain.hostname
  }

  // convert unix timestamp to yyyy-mm-dd
  const formatTimestamp = (v) => {
    let date = new Date(v.seconds * 1000)
    var day = date.getDate();
    // month is from 0 to 11
    var month = date.getMonth() + 1;
    var year = date.getFullYear();
    return '' + year + '-' + (month <= 9 ? '0' + month : month) + '-' + (day <= 9 ? '0' + day : day);
  }

  return (
    <Row className='bg-white shadow-sm row-cols-1 row-cols-md-2 my-2'>
      <Col>
        <div className='d-flex align-items-center'>
          <div className=''>
            <img src={"https://www.google.com/s2/favicons?sz=32&domain_url=" + domainName(props.data.url)} loading="lazy" alt="logo" />
          </div>
          <div className='py-2 text-left mx-2'>
            <GradientText className="fw-bold text-truncate"><a href={props.data.short_url} target="_blank"
              rel="noreferrer" title={"Shortened URL for " + props.data.url}>
              {(props.data.title !== null && props.data.title !== "") ? props.data.title : props.data.short_url}</a>
            </GradientText>
            <div className="small text-truncate">{props.data.url}</div>
            <div className='d-flex'>
              <div className="small text-truncate me-2">{formatTimestamp(props.data.created_at)}</div>
              {props.data.tags.slice(0, 3).map((item, key) => (
                <Badge key={key} bg="secondary" className="small text-truncate me-1">{item}</Badge>
              ))}
              {props.data.tags.length > 2 &&
                <OverlayTrigger placement="top" overlay={<Tooltip id="clicks">{props.data.tags.toString()}</Tooltip>}>
                  <Badge bg="secondary" className="small text-truncate me-1">...</Badge>
                </OverlayTrigger>
              }
            </div>
          </div>
        </div>
      </Col>
      <Col>
        <div className='py-3'>
          <OverlayTrigger placement="top" overlay={<Tooltip id="clicks">Copy</Tooltip>}>
            <Button className="me-1 my-1 float-end" variant='outline-primary' data-hash={props.data.short_url} onClick={(el) => props.copy(el)}>
              <RiFileCopyLine />
            </Button>
          </OverlayTrigger>
          <OverlayTrigger placement="top" overlay={<Tooltip id="clicks">Delete</Tooltip>}>
            <Button className="me-1 my-1 float-end" variant='outline-primary' data-hash={props.data.alias} onClick={(el) => props.delete(el)}>
              <RiDeleteBinLine />
            </Button>
          </OverlayTrigger>
          <OverlayTrigger placement="top" overlay={<Tooltip id="clicks">Edit</Tooltip>}>
            <Button className="me-1 my-1 float-end" variant='outline-primary' data-hash={props.data.alias} onClick={(el) => props.edit(el)}>
              <FiEdit />
            </Button>
          </OverlayTrigger>
          <OverlayTrigger placement="top" overlay={<Tooltip id="clicks">Clicks</Tooltip>}>
            <Button className="me-1 my-1 float-end" variant='outline-primary' data-hash={props.data.alias} onClick={(el) => props.showStat(el)}>
              <BiBarChart />
            </Button>
          </OverlayTrigger>
          <OverlayTrigger placement="top" overlay={<Tooltip id="clicks">QR code</Tooltip>}>
            <Button className="me-1 my-1 float-end" variant='outline-primary' data-hash={props.data.alias} onClick={(el) => props.qrcode(el)}>
              <FaQrcode />
            </Button>
          </OverlayTrigger>
          <Dropdown>
            <OverlayTrigger placement="top" overlay={<Tooltip id="clicks">Share</Tooltip>}>
              <Dropdown.Toggle className="me-1 my-1 float-end" variant="primary" id="dropdown-basic">
                <RiShareBoxLine />
              </Dropdown.Toggle>
            </OverlayTrigger>
            <Dropdown.Menu>
              <Dropdown.Item>
                <FacebookShareButton url={props.data.short_url} title="Share ShortUrl on Facebook" description="" quote="Share ShortUrl on Facebook" hashtag="">
                  <AiOutlineFacebook className="me-2" />
                  Facebook
                </FacebookShareButton>
              </Dropdown.Item>
              <Dropdown.Item>
                <TwitterShareButton url={props.data.short_url} title="Share ShortUrl on Twitter">
                  <RiTwitterLine className="me-2" />
                  Twitter
                </TwitterShareButton>
              </Dropdown.Item>
              <Dropdown.Item>
                <WhatsappShareButton url={props.data.short_url} title="Share ShortUrl on WhatsApp" separator=":: ">
                  <AiOutlineWhatsApp className="me-2" />
                  WhatsApp
                </WhatsappShareButton>
              </Dropdown.Item>
              <Dropdown.Item>
                <LinkedinShareButton url={props.data.short_url}>
                  <AiOutlineLinkedin className="me-2" />
                  LinkedIn
                </LinkedinShareButton>
              </Dropdown.Item>
              <Dropdown.Item>
                <EmailShareButton url={props.data.short_url} subject="Share ShortUrl via Email">
                  <AiOutlineMail className="me-2" />
                  Email
                </EmailShareButton>
              </Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        </div>
      </Col>
    </Row>
  )
}
export default ShortUrlItem;