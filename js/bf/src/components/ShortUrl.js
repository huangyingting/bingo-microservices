import React, { useState, useEffect } from 'react'
import { Col, Row, Form, Button, FloatingLabel } from 'react-bootstrap';
import "react-circular-progressbar/dist/styles.css";
import ShortUrlItem from './ShortUrlItem';
import useFetch from 'use-http'
import { API_ENDPOINT } from '../Global';
import { GetAccessToken } from '../AAD'
import QRCode from 'qrcode.react';
import swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content'
import { Section } from './Styled';
import ShortUrlEdit from "./ShortUrlEdit";
import ShortUrlQnA from './ShortUrlQnA';
import ShortUrlToast from "./ShortUrlToast";
import ShortUrlHeader from './ShortUrlHeader';
import { useForm } from "react-hook-form";
import { PAGE_SIZE } from "../Global"
import { MdOutlineNavigateBefore, MdOutlineNavigateNext } from "react-icons/md";

function compareArray(array1, array2) {
  const array2Sorted = array2.slice().sort();
  return array1.length === array2.length && array1.slice().sort().every(function (value, index) {
    return value === array2Sorted[index];
  });
}

function ShortUrl() {
  // short url list
  const [shortUrls, setShortUrls] = useState([]);
  // current short url
  const [currentShortUrl, setCurrentShortUrl] = useState(null)
  // error message
  const [errorMsg, setErrorMsg] = useState("")
  // copied
  const [infoMsg, setInfoMsg] = useState("")
  // show edit
  const [editVisible, setEditVisible] = useState(false)
  // dialog
  const swalReact = withReactContent(swal)
  // refresh
  const [refresh, setRefresh] = useState({ start: 0, count: PAGE_SIZE })
  // pagination
  const [pagination, setPagination] = useState({ start: 0, count: 0 })

  const { register, handleSubmit, reset, formState: { errors } } = useForm({
    criteriaMode: "all",
  });

  // Use fetch
  const { response, del, get, post, put } = useFetch(API_ENDPOINT, {
    retries: 0,
    cachePolicy: "no-cache",
    interceptors: {
      request: async ({ options }) => {
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

  const showErrorMsg = (msg) => {
    setErrorMsg(msg)
  }

  const hideErrorMsg = () => {
    setErrorMsg("")
  }

  const showInfoMsg = (msg) => {
    setInfoMsg(msg)
  }

  const hideInfoMsg = () => {
    setInfoMsg("")
  }

  // Create short url
  const createShortUrl = async (v) => {
    var r = await post('/v1/shorturl', { "url": v.url, "alias": v.alias })
    if (response.ok) {
      /*
      setShortUrls([{
        "url": r.url,
        "alias": r.alias,
        "short_url": API_ENDPOINT + "/" + r.alias,
        "title": r.title,
        "tags": r.tags,
        "fraud_detection": r.fraud_detection,
        "disabled": r.disabled,
        "no_referrer": r.no_referrer,
        "utm_source": r.utm_source,
        "utm_medium": r.utm_medium,
        "utm_campaign": r.utm_campaign,
        "utm_term": r.utm_term,
        "utm_content": r.utm_content,
        "created_at": r.created_at
      }, ...shortUrls.slice(0, 9)])
      */
      setRefresh({ start: 0, count: PAGE_SIZE })
      showInfoMsg("Short url " + API_ENDPOINT + "/" + r.alias + " is created")
      reset({ url: "", alias: "" })
    } else {
      showErrorMsg(r.message)
    }
  }

  // Delete short url
  const deleteShortUrl = async (el) => {
    el.preventDefault()
    const alias = el.currentTarget.getAttribute('data-hash');

    swalReact.fire({
      title: 'Are you sure?',
      text: "You won't be able to revert this!",
      icon: 'warning',
      showCancelButton: true,
      cancelButtonColor: '#d33',
      confirmButtonText: 'Yes, delete it!',
      confirmButtonColor: '#4da6e7',
    }).then(async (result) => {
      if (result.isConfirmed) {
        try {
          var r = await del('/v1/shorturl/' + alias)
          if (response.ok) {
            // trigger useEffect to list short url again
            setRefresh({ start: pagination.start, count: PAGE_SIZE })
            showInfoMsg("Alias " + alias + " is deleted")
          } else {
            showErrorMsg(r.message)
          }
        } catch (err) {
          throw new Error(err);
        }
      }
    })
  }

  // Edit short url
  const editShortUrl = (el) => {
    el.preventDefault()
    const alias = el.currentTarget.getAttribute('data-hash');
    const shortUrl = shortUrls.find(x => x.alias === alias)
    setCurrentShortUrl(shortUrl)
    setEditVisible(true)
  }

  const updateShortUrl = async (v) => {
    const p_shortUrl = shortUrls.find(x => x.alias === v.alias)
    if (p_shortUrl.url !== v.url ||
      p_shortUrl.title !== v.title ||
      !compareArray(p_shortUrl.tags, v.tags) ||
      p_shortUrl.fraud_detection !== v.fraud_detection ||
      p_shortUrl.disabled !== v.disabled ||
      p_shortUrl.no_referrer !== v.no_referrer ||
      p_shortUrl.utm_source !== v.utm_source ||
      p_shortUrl.utm_medium !== v.utm_medium ||
      p_shortUrl.utm_campaign !== v.utm_campaign ||
      p_shortUrl.utm_term !== v.utm_term ||
      p_shortUrl.utm_content !== v.utm_content) {
      var r = await put('/v1/shorturl', {
        "alias": v.alias, "url": v.url,
        "title": v.title, "tags": v.tags,
        "fraud_detection": v.fraud_detection,
        "disabled": v.disabled,
        "no_referrer": v.no_referrer,
        "utm_source": v.utm_source,
        "utm_medium": v.utm_medium,
        "utm_campaign": v.utm_campaign,
        "utm_term": v.utm_term,
        "utm_content": v.utm_content
      })
      if (response.ok) {
        setShortUrls(
          shortUrls.map(item => {
            return item.alias === r.alias
              ? {
                ...item,
                url: r.url,
                title: r.title,
                tags: r.tags,
                fraud_detection: r.fraud_detection,
                disabled: r.disabled,
                no_referrer: r.no_referrer,
                utm_source: r.utm_source,
                utm_medium: r.utm_medium,
                utm_campaign: r.utm_campaign,
                utm_term: r.utm_term,
                utm_content: r.utm_content
              }
              : item
          }
          ))
        showInfoMsg("Short url " + API_ENDPOINT + "/" + r.alias + " is updated")
      } else {
        showErrorMsg(r.message)
      }
    } else {
      showInfoMsg("Won't update as nothing changed")
    }
  }


  const suggestedTags = async (input, callback) => {
    const createOption = (label) => ({
      label,
      value: label,
    });
    const r = await get('/v1/tag-suggest/' + input)
    if (response.ok) {
      callback(r.value.map(tag => {
        return createOption(tag)
      }))
    }
    // return an empty array if anything goes wrongly
    return []
  }


  // Copy short url to clipboard
  const copyShortUrl = (el) => {
    el.preventDefault()
    const short_url = el.currentTarget.getAttribute('data-hash');
    navigator.clipboard.writeText(short_url)
    showInfoMsg("Short url " + short_url + " is copied to clipboard")
  }

  const downloadQRCode = (alias) => {
    const canvas = document.getElementById(alias);
    const pngUrl = canvas
      .toDataURL("image/png")
      .replace("image/png", "image/octet-stream");
    let downloadLink = document.createElement("a");
    downloadLink.href = pngUrl;
    downloadLink.download = alias + ".png";
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
    showInfoMsg("QRCode image downloaded")
  }

  // Generate QR code
  const generateQRCode = (el) => {
    el.preventDefault()
    const alias = el.currentTarget.getAttribute('data-hash');
    const shortUrl = shortUrls.find(x => x.alias === alias)
    swalReact.fire({
      title: "QRCode",
      html: (
        <>
          <div>Click image to download QRCode</div>
          <QRCode size={240} id={alias} value={shortUrl.short_url} onClick={() => downloadQRCode(alias)} />
        </>
      ),
      confirmButtonText: "Close",
      confirmButtonColor: '#4da6e7',
    })
  }

  // count clicks
  const countClicks = async (el) => {
    el.preventDefault()
    const alias = el.currentTarget.getAttribute('data-hash');
    const r = await get('/v1/shorturl-bi/clicks/' + alias)
    if (response.ok) {
      swalReact.fire({
        title: "Number of clicks",
        html: (<h3>{r.clicks}</h3>),
        confirmButtonColor: '#4da6e7'
      })
    } else {
      showErrorMsg(r.message)
    }
  }

  // load all short urls
  useEffect(() => {
    const listShortUrls = async () => {
      const shortUrls = await get('/v1/shorturl?start=' + refresh.start + '&count=' + refresh.count)
      if (response.ok) {
        if (shortUrls.value != null) {
          setShortUrls(shortUrls.value.map(i => {
            return {
              "url": i.url,
              "alias": i.alias,
              "short_url": API_ENDPOINT + "/" + i.alias,
              "title": i.title,
              "tags": i.tags,
              "fraud_detection": i.fraud_detection,
              "disabled": i.disabled,
              "no_referrer": i.no_referrer,
              "utm_source": i.utm_source,
              "utm_medium": i.utm_medium,
              "utm_campaign": i.utm_campaign,
              "utm_term": i.utm_term,
              "utm_content": i.utm_content,
              "created_at": i.created_at
            }
          }).slice(0, PAGE_SIZE))
        }
        setPagination({ start: shortUrls.start, count: shortUrls.count })
      }
    }
    listShortUrls()
  }, [refresh]);

  const nextPage = () => {
    setRefresh({ start: pagination.start + PAGE_SIZE, count: PAGE_SIZE })
  }

  const prevPage = () => {
    setRefresh({ start: pagination.start - PAGE_SIZE, count: PAGE_SIZE })
  }

  return (
    <>
      <Section color="#e1f4fd">
        <ShortUrlHeader />
        <Row className="justify-content-center align-items-center pt-2 pb-2">
          <Col className="col-12 col-md-11 col-lg-10 col-xl-9 col-xxl-8">
            <Form onSubmit={handleSubmit((data) => { createShortUrl(data) })}>
              <Row className="justify-content-center pb-4 row-cols-1 row-cols-md-2">
                <Col className="col-12 col-md-8 px-0">
                  <FloatingLabel controlId="floatingUrl" label={errors.url ? errors.url?.message : "Shorten your link"}>
                    <Form.Control autoComplete="off" isInvalid={errors.url} type="text" placeholder="Shorten your Link" className="fw-bold" {
                      ...register("url", {
                        required: 'Url is required',
                        pattern: {
                          value: /^(?:(?:(?:https?|http):)?\/\/)(?:\S+(?::\S*)?@)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:[/?#]\S*)?$/i,
                          message: "Invalid URL",
                        }
                      })} />
                  </FloatingLabel>
                </Col>
                <Col className="col-12 col-md-4 px-0">
                  <div className="d-flex">
                    <FloatingLabel className="flex-fill" controlId="floatingAlias" label={errors.alias ? errors.alias?.message : "Alias"}>
                      <Form.Control autoComplete="off" isInvalid={errors.alias} type="text" placeholder="Alias" className="fw-bold" {
                        ...register("alias", {
                          pattern: {
                            value: /^[0-9A-Za-z_-]*$/i,
                            message: 'Invalid alias',
                          },
                        })} />
                    </FloatingLabel>
                    <Button variant='primary' type="submit">Shorten</Button>
                  </div>
                </Col>
              </Row>
            </Form>
          </Col>
        </Row>
        <Row className="justify-content-center align-items-center pt-2 pb-2">
          <Col className="col-12 col-md-11 col-lg-10 col-xl-9 col-xxl-8">
            {shortUrls.map((item, key) => (
              <ShortUrlItem key={item.alias} data={item} copy={copyShortUrl} qrcode={generateQRCode}
                delete={deleteShortUrl} edit={editShortUrl} showStat={countClicks} />
            ))}
            <Row className='row-cols-1 row-cols-md-2'>
              <Col></Col>
              <Col>
                <Button className="float-end" variant='primary' disabled={pagination.count < PAGE_SIZE} onClick={nextPage}>
                  <MdOutlineNavigateNext />
                </Button>
                <Button className="me-1 float-end" variant='primary' disabled={pagination.start === 0} onClick={prevPage}>
                  <MdOutlineNavigateBefore />
                </Button>
              </Col>
            </Row>
          </Col>
        </Row>


      </Section>
      {currentShortUrl !== null &&
        <ShortUrlEdit visible={editVisible} hide={() => setEditVisible(false)}
          update={updateShortUrl} loadOptions={suggestedTags} data={currentShortUrl} />
      }
      <ShortUrlQnA />
      <ShortUrlToast info={infoMsg} error={errorMsg} hideInfo={hideInfoMsg} hideError={hideErrorMsg} />
    </>
  )
}
export default ShortUrl;