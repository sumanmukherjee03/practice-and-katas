const Input = (props) => {
  return (
    <div className="mb-3">
      {/* Notice the use of htmlFor instead of "for" in JSX. Also note the use of JS templating in the input element className attribute. */}
      <label htmlFor={props.name} className="form-label">{props.title}</label>
      <input type={props.type} className={`form-control ${props.className}`} id={props.name} name={props.name} value={props.value} onChange={props.handleChange} placeholder={props.placeholder} />
      <div className={props.errorDiv}>{props.errorMsg}</div>
    </div>
  );
}
export default Input;
