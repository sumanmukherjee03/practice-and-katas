const Select = (props) => {
  return (
    <div className="mb-3">
      {/* Notice the use of htmlFor instead of "for" in JSX */}
      <label htmlFor={props.name} className="form-label">{props.title}</label>
      <select className="form-select" id={props.name} name={props.name} value={props.value} onChange={props.handleChange} >
        <option className="form-select" value="">{props.placeholder}</option>
        {Object.keys(props.options).map((optKey) => (
          <option className="form-select" key={optKey} value={optKey}>{props.options[optKey]}</option>
        ))}
      </select>
    </div>
  );
}
export default Select;
