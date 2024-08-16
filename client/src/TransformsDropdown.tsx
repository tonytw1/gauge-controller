export function TransformsDropdown({transforms}: { transforms: Transform[] }) {

    function TransformOption({transform}: { transform: Transform }) {
        return (
            <>
                <option value={transform.Name}>{transform.Name}</option>
            </>
        )
    }

    const listItems = transforms.map(transform => <TransformOption transform={transform}/>);

    return (
        <>
            <select name="Transform">{listItems}</select>
        </>
    )
}
