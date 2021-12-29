const TextArea = (props) => {
  return (
    <div className="mb-3">
      {/* Notice the use of htmlFor instead of "for" in JSX */}
      <label htmlFor={props.name} className="form-label">{props.title}</label>
      <textarea className="form-control" id={props.name} name={props.name} rows="3" onChange={props.handleChange} value={props.value} placeholder={props.placeholder} />
    </div>
  );
}
export default TextArea;
