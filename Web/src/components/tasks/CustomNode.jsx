import { Handle, Position } from 'reactflow'
import './CustomNode.css'

function CustomNode({ data }) {
  return (
    <div className="custom-node">
      <Handle type="target" position={Position.Left} />
      <div className="node-content">
        <span className="node-label">{data.label}</span>
      </div>
      <Handle type="source" position={Position.Right} />
    </div>
  )
}

export default CustomNode
