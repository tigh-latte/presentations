local M = {
	state = { "> " },
	is_markdown = false,
}

function M.plugin(input)
	M.state = { input.initial or "> " }

	return M.state[1] .. "\n"
end

---@param input string
function M.on_key(input)
	if input == "backspace" then
		if #M.state > 1 then
			table.remove(M.state, #M.state)
		end
	elseif input == "enter" then
	else
		table.insert(M.state, input)
	end

	local line = table.concat(M.state, "")

	return line .. "\n"
end

return M
