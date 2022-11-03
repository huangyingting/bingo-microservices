import { Form, Button, Offcanvas, FloatingLabel, InputGroup } from 'react-bootstrap';
import { useForm, Controller } from "react-hook-form";
import React, { useMemo, useEffect } from 'react'
import { GradientText } from './Styled';
import AsyncCreatableSelect from 'react-select/async-creatable';

const ShortUrlEdit = (props) => {
  const { register, getValues, handleSubmit, control, reset, formState: { errors } } = useForm({
    defaultValues: useMemo(() => {
      return {
        tags: props.data.tags,
        url: props.data.url,
        alias: props.data.alias,
        title: props.data.title,
        fraud_detection: props.data.fraud_detection,
        disabled: props.data.disabled,
        no_referrer: props.data.no_referrer,
        utm_source: props.data.utm_source,
        utm_medium: props.data.utm_medium,
        utm_campaign: props.data.utm_campaign,
        utm_term: props.data.utm_term,
        utm_content: props.data.utm_content
      };
    }, [props]),
    mode: "onChange",
    criteriaMode: "all",
  });

  useEffect(() => {
    reset({
      alias: props.data.alias,
      url: props.data.url,
      title: props.data.title,
      tags: props.data.tags,
      fraud_detection: props.data.fraud_detection,
      disabled: props.data.disabled,
      no_referrer: props.data.no_referrer,
      utm_source: props.data.utm_source,
      utm_medium: props.data.utm_medium,
      utm_campaign: props.data.utm_campaign,
      utm_term: props.data.utm_term,
      utm_content: props.data.utm_content
    });
  }, [reset, props.data]);

  const updateShortUrl = v => {
    props.hide()
    console.log(v)
    props.update(v)
  }

  function isEmpty(str) {
    return (!str || str.length === 0)
  }

  return (
    <Offcanvas show={props.visible} onHide={props.hide} placement="end">
      <Offcanvas.Header closeButton className='border-bottom'>
        <Form.Label className="fw-bold"><GradientText>{props.data.short_url}</GradientText></Form.Label>
      </Offcanvas.Header>
      <Offcanvas.Body>
        <Form onSubmit={handleSubmit((data) => { updateShortUrl(data) })}>
          <Form.Group>
            <Form.Label className="fw-bold">Basic</Form.Label>
            <FloatingLabel controlId="floatingTitle" label={errors.title ? errors.title?.message : "Url Title"}>
              <Form.Control autoComplete="off" isInvalid={errors.title} type="text" placeholder="Url Title" {
                ...register("title", {
                  pattern: {
                    value: /^[0-9A-Za-z_-]*$/i,
                    message: 'Invalid title',
                  },
                })} />
            </FloatingLabel>
          </Form.Group>
          <Form.Group className="mb-1">
            <FloatingLabel controlId="floatingUrl" label={errors.url ? errors.url?.message : "Url"}>
              <Form.Control autoComplete="off" isInvalid={errors.url} type="text" placeholder="Url is required" {
                ...register("url", {
                  required: 'Url is required',
                  pattern: {
                    value: /^(?:(?:(?:https?|http):)?\/\/)(?:\S+(?::\S*)?@)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:[/?#]\S*)?$/i,
                    message: "Invalid URL",
                  }
                })} />
            </FloatingLabel>
          </Form.Group>
          <Form.Group className="mb-1">
            <Form.Label className="fw-bold">Tags</Form.Label>
            <Controller
              name="tags"
              control={control}
              rules={{
                validate: tags => {
                  var pattern = /^[0-9A-Za-z_-]*$/i
                  for (const tag of tags) {
                    if (!pattern.test(tag)) {
                      return "Invalid tag: " + tag
                    }
                  }
                  return true
                }
              }}
              render={({ field }) =>
                <AsyncCreatableSelect
                  value={field.value.map(v => { return { value: v, label: v } })}
                  cacheOptions
                  loadOptions={props.loadOptions}
                  isClearable
                  isMulti
                  placeholder="Type and press enter"
                  onChange={(e) => field.onChange(e.map(i => { return i.value }))} />}
            />
            <small className='text-danger'>{errors.tags?.message}</small>
          </Form.Group>
          <Form.Group className="mb-1">
            <Form.Label className="fw-bold">Flags</Form.Label>
            <Controller
              name="fraud_detection"
              control={control}
              render={({ field }) =>
                <Form.Check
                  type="switch"
                  checked={field.value}
                  label="Enable Fraud Detection"
                  onChange={(e) => field.onChange(e.target.checked ? true : false)} />}
            />
            <Controller
              name="no_referrer"
              control={control}
              render={({ field }) =>
                <Form.Check
                  type="switch"
                  checked={field.value}
                  label="Hide referrer"
                  onChange={(e) => field.onChange(e.target.checked ? true : false)} />}
            />
            <Controller
              name="disabled"
              control={control}
              render={({ field }) =>
                <Form.Check
                  type="switch"
                  checked={field.value}
                  label="Disable short link"
                  onChange={(e) => field.onChange(e.target.checked ? true : false)} />}
            />
          </Form.Group>
          <Form.Group className="mb-1">
            <Form.Label className="fw-bold">Campaign Tracking</Form.Label>
            <InputGroup size="sm">
              <InputGroup.Text className="w-25 fw-bold">Source</InputGroup.Text>
              <Form.Control autoComplete="off" isInvalid={errors.utm_source} type="text" placeholder="e.g. twitter, facebook" {
                ...register("utm_source", {
                  deps: ['utm_medium', 'utm_campaign'],
                  validate: {
                    required: value => {
                      if (isEmpty(value) && (!isEmpty(getValues("utm_medium")) ||
                        !isEmpty(getValues("utm_campaign")) ||
                        !isEmpty(getValues("utm_term")) ||
                        !isEmpty(getValues("utm_content")))) {
                        return 'Source field is required';
                      }
                      return true;
                    },
                  },
                  pattern: {
                    value: /^[0-9A-Za-z_-]*$/i,
                    message: 'Invalid characters',
                  },
                })} />
            </InputGroup>
            <small className='text-danger'>
              {errors.utm_source?.message}
            </small>
            <InputGroup size="sm">
              <InputGroup.Text className="w-25 fw-bold">Medium</InputGroup.Text>
              <Form.Control autoComplete="off" isInvalid={errors.utm_medium} type="text" placeholder="e.g. banner, email" {
                ...register("utm_medium", {
                  deps: ['utm_source', 'utm_campaign'],
                  validate: {
                    required: value => {
                      if (isEmpty(value) && (!isEmpty(getValues("utm_source")) ||
                        !isEmpty(getValues("utm_campaign")) ||
                        !isEmpty(getValues("utm_term")) ||
                        !isEmpty(getValues("utm_content")))) {
                        return 'Medium field is required';
                      }
                      return true;
                    },
                  },
                  pattern: {
                    value: /^[0-9A-Za-z_-]*$/i,
                    message: 'Invalid characters',
                  },
                })} />
            </InputGroup>
            <small className='text-danger'>
              {errors.utm_medium?.message}
            </small>
            <InputGroup size="sm">
              <InputGroup.Text className="w-25 fw-bold">Campaign</InputGroup.Text>
              <Form.Control autoComplete="off" isInvalid={errors.utm_campaign} type="text" placeholder="e.g. spring_sales" {
                ...register("utm_campaign", {
                  deps: ['utm_source', 'utm_medium'],
                  validate: {
                    required: value => {
                      if (isEmpty(value) && (!isEmpty(getValues("utm_source")) ||
                        !isEmpty(getValues("utm_medium")) ||
                        !isEmpty(getValues("utm_term")) ||
                        !isEmpty(getValues("utm_content")))) {
                        return 'Campaign field is required';
                      }
                      return true;
                    },
                  },
                  pattern: {
                    value: /^[0-9A-Za-z_-]*$/i,
                    message: 'Invalid characters',
                  },
                })} />
            </InputGroup>
            <small className='text-danger'>
              {errors.utm_campaign?.message}
            </small>
            <InputGroup size="sm">
              <InputGroup.Text className="w-25">Term</InputGroup.Text>
              <Form.Control autoComplete="off" isInvalid={errors.utm_term} type="text" placeholder="e.g. sales, shoes" {
                ...register("utm_term", {
                  deps: ['utm_source', 'utm_medium', 'utm_campaign'],
                  pattern: {
                    value: /^[0-9A-Za-z_-]*$/i,
                    message: 'Invalid characters',
                  },
                })} />
            </InputGroup>
            <small className='text-danger'>
              {errors.utm_term?.message}
            </small>
            <InputGroup size="sm">
              <InputGroup.Text className="w-25">Content</InputGroup.Text>
              <Form.Control autoComplete="off" isInvalid={errors.utm_content} type="text" placeholder="e.g. logolink, textlink" {
                ...register("utm_content", {
                  deps: ['utm_source', 'utm_medium', 'utm_campaign'],
                  pattern: {
                    value: /^[0-9A-Za-z_-]*$/i,
                    message: 'Invalid characters',
                  },
                })} />
            </InputGroup>
            <small className='text-danger'>
              {errors.utm_content?.message}
            </small>
          </Form.Group>
          <Button variant="primary" type="submit">
            Update
          </Button>
        </Form>
      </Offcanvas.Body>
    </Offcanvas>
  )
}

export default ShortUrlEdit;