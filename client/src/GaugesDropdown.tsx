import {Select} from '@headlessui/react'

export function GaugesDropdown({gauges}: { gauges: Gauge[] }) {
    return (
        <Select name="Gauge" aria-label="Destination gauge">
            {
                gauges.map(gauge =>
                    <option value={gauge.Name}>{gauge.Name}</option>
                )
            }
        </Select>
    )
}
