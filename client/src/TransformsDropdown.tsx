import {Select} from '@headlessui/react'

export function TransformsDropdown({transforms}: { transforms: Transform[] }) {
    return (
        <Select name={"Transform"} aria-description={"Transform to apply"}>
            {transforms.map(transform =>
                <option value={transform.Name}>{transform.Name}</option>
            )
            }
        </Select>
    )
}
